package models

import "time"

type RSA struct {
	ID               int       `gorm:"primaryKey" json:"id,omitempty"`
	UserID           uint      `gorm:"column:user_id;not null" json:"user_id,omitempty"`
	EncryptedMessage string    `gorm:"column:encrypted_message;type:text;not null" json:"encrypted_message,omitempty"`
	D                string    `gorm:"column:d;type:text;not null" json:"d,omitempty"`
	N                string    `gorm:"column:n;type:text;not null" json:"n,omitempty"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
}
