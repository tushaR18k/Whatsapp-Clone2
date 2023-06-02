package main

import (
	"fmt"
	"net/http"
	"whatsapp-backend/middleware"
	"whatsapp-backend/models"
	"whatsapp-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func main() {
	port := "8080"

	db := routes.GetDbConnection()

	//Auto Migrating the user Model
	err := MigrateTables(db)
	if err != nil {
		fmt.Println("Error in db: ", err)
	}

	manager := routes.NewManager()

	//creating a new router
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type"}
	router.Use(cors.New(corsConfig))

	//Apply the AuthMiddleware to specific routes
	// authRoutes := router.Group("/api")
	// authRoutes.Use(middleware.AuthMiddleware())

	router.POST("/api/signup", routes.Signup)
	router.POST("/api/login", routes.Login)
	router.POST("/api/add-friend", middleware.AuthMiddleware(), routes.AddFriend)
	router.GET("/api/friends/:userId", middleware.AuthMiddleware(), routes.GetFriends)

	router.POST("/api/messages", middleware.AuthMiddleware(), routes.SendMessage)
	router.GET("/api/messages/:senderID/:receiverID", middleware.AuthMiddleware(), routes.GetMessages)

	router.GET("/ws", manager.ServeWS)
	router.GET("/api/files/:filePath", routes.DownloadFile)

	router.Run(":" + port)

}

func MigrateTables(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{}, &models.Friend{}, &models.Message{}).Error
	if err != nil {
		return err
	}

	err = db.Model(&models.Friend{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		return err
	}

	err = db.Model(&models.Message{}).AddForeignKey("sender_id", "users(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		return err
	}

	err = db.Model(&models.Message{}).AddForeignKey("receiver_id", "users(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		return err
	}

	return nil

}
