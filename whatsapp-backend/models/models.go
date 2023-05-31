package models

import (
	"time"
)

type User struct {
	ID             int        `gorm:"primary_key" json:"id"`
	Username       string     `gorm:"unique:not null" json:"username"`
	Email          string     `gorm:"unique not null" json:"email"`
	Password       string     `gorm:"not null" json:"password"`
	Token          string     `gorm:"-" json:"token"`
	TokenCreatedAt *time.Time `gorm:"column:token_created_at" json:"token_created_at"`
	TokenExpiresAt *time.Time `gorm:"column:token_expires_at" json:"token_expires_at"`
	Friends        []Friend   `gorm:"foreignKey:UserID"`
}

type Friend struct {
	ID        int       `gorm:"primary_key" json:"id"`
	UserID    int       `gorm:"not null" json:"user_id"`
	FriendId  int       `gorm:"not null" json:"friend_id"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
}

type Message struct {
	ID          int       `gorm:"primary_key" json:"id"`
	SenderID    int       `gorm:"not null" json:"sender_id"`
	ReceiverID  int       `gorm:"not null" json:"receiver_id"`
	MessageType string    `gorm:"not null" json:"message_type"`
	Content     string    `gorm:"type:text" json:"context"`
	FileName    string    `json:"file_name"`
	FileSize    string    `json:"file_size"`
	FileType    string    `json:"file_type"`
	FilePath    string    `json:"file_path"`
	CreatedAt   time.Time `gorm:"default:current_timestamp"`

	Sender   User `gorm:"foreignkey:SenderID" json:"-"`
	Receiver User `gorm:"foreignkey:ReceiverID" json:"-"`
}

type WebSocketResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
