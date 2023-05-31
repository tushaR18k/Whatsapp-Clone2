package routes

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"whatsapp-backend/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	secret = "my_secret"
)

type AddFriendRequest struct {
	UserID      int    `json:"userId"`
	FriendEmail string `json:"friendEmail"`
}

type UserID struct {
	ID int `json:"userId"`
}

type FriendList struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func Signup(c *gin.Context) {
	var user models.User
	db := GetDbConnection()
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().UTC().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	now := time.Now()
	user.Token = tokenString
	user.TokenCreatedAt = &now
	expTime := time.Unix(token.Claims.(jwt.MapClaims)["exp"].(int64), 0)
	user.TokenExpiresAt = &expTime

	//Save the user to the database
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	c.JSON(http.StatusOK, tokenString)

}

func Login(c *gin.Context) {
	var user models.User
	db := GetDbConnection()
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var dbUser models.User
	if err := db.Where("email=?", user.Email).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	//comparing the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": dbUser.Username,
		"email":    dbUser.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	now := time.Now()
	dbUser.Token = tokenString
	dbUser.TokenCreatedAt = &now
	expTime := time.Unix(token.Claims.(jwt.MapClaims)["exp"].(int64), 0)
	user.TokenExpiresAt = &expTime

	if err := db.Save(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "id": dbUser.ID})

}

func AddFriend(c *gin.Context) {
	db := GetDbConnection()
	var req AddFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var currentUser models.User
	if err := db.First(&currentUser, req.UserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	if currentUser.Email == req.FriendEmail {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot add your ownself as a friend"})
		return
	}

	//Query db to see friend exists
	var friendUser models.User
	if err := db.Where("email=?", req.FriendEmail).First(&friendUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Friend's email ID does not exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check the emailID"})
		return
	}

	isFrnd := IsFriend(req.UserID, friendUser.ID)
	if !isFrnd {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You are already friends."})
		return
	}

	//Friend does exist
	friend := models.Friend{
		UserID:   req.UserID,
		FriendId: friendUser.ID,
	}

	// Add the user as a friend for the friend added by the user
	userFriend := models.Friend{
		UserID:   friendUser.ID,
		FriendId: req.UserID,
	}

	if err := db.Create(&friend).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add friend"})
		return
	}

	if err := db.Create(&userFriend).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to user as a friend"})
		return
	}

	currentUser.Friends = append(currentUser.Friends, friend)
	if err := db.Save(&currentUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user's friend list"})
		return
	}

	friendUser.Friends = append(friendUser.Friends, userFriend)
	if err := db.Save(&currentUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update friend's friend list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend added successfully"})
}

func GetFriends(c *gin.Context) {

	userIdStr := c.Param("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	db := GetDbConnection()
	var friends []FriendList
	result := db.Table("friends").
		Select("users.id, users.email").
		Joins("JOIN users ON users.id = friends.friend_id").
		Where("friends.user_id = ?", userId).
		Find(&friends)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching friends list"})
		return
	}
	fmt.Print(friends)

	c.JSON(http.StatusOK, friends)

}

func GetDbConnection() *gorm.DB {
	dbString := "postgres://postgres:tushar@localhost:5432/whatsapp?sslmode=disable"
	db, err := gorm.Open("postgres", dbString)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func IsFriend(userID int, friendId int) bool {

	db := GetDbConnection()
	var friend models.Friend
	result := db.Joins("JOIN users ON users.id = friends.user_id").Where("users.id = ? AND friends.friend_id = ?", userID, friendId).First(&friend)
	return result.Error == gorm.ErrRecordNotFound

}

func SendMessage(c *gin.Context) {
	//Parsing the body
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate required fields
	if message.SenderID == 0 || message.ReceiverID == 0 || message.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	//Saving the message
	db := GetDbConnection()
	if err := db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message Sent Successfully"})
}

func GetMessages(c *gin.Context) {
	//Parsing the sender and receiver ID
	senderIDStr := c.Param("senderID")
	receiverIDStr := c.Param("receiverID")
	senderID, _ := strconv.Atoi(senderIDStr)
	receiverID, _ := strconv.Atoi(receiverIDStr)

	//Parsing the page number from query param
	pageNumberStr := c.Query("page")
	pageNumber, _ := strconv.Atoi(pageNumberStr)
	pageSize := 10

	//Calculating the offset based on the page number and page size
	offset := (pageNumber - 1) * pageSize

	//Retrieving the messages from the database
	db := GetDbConnection()
	var messages []models.Message
	if err := db.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		senderID, receiverID, receiverID, senderID).Offset(offset).Limit(pageSize).Order("created_at").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
