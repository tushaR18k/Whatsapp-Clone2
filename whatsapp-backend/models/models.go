package models

import (
	"mime/multipart"
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
	ID          int                   `gorm:"primary_key" json:"id"`
	SenderID    int                   `gorm:"not null" json:"sender_id" form:"sender_id"`
	ReceiverID  int                   `gorm:"not null" json:"receiver_id" form:"receiver_id"`
	MessageType string                `gorm:"not null" json:"message_type" form:"message_type"`
	Content     string                `gorm:"type:text" json:"context" form:"context"`
	FileName    string                `gorm:"column:file_name" json:"file_name,omitempty"`
	FileSize    string                `gorm:"column:file_size" json:"file_size,omitempty"`
	FileType    string                `gorm:"column:file_type" json:"file_type,omitempty"`
	FilePath    string                `gorm:"column:file_path" json:"file_path,omitempty"`
	FileContent *multipart.FileHeader `form:"file"`
	CreatedAt   time.Time             `gorm:"default:current_timestamp"`

	Sender   User `gorm:"foreignkey:SenderID" json:"-"`
	Receiver User `gorm:"foreignkey:ReceiverID" json:"-"`
}

type WebSocketResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
