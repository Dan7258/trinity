package models

import "time"

type Kuznechik struct {
	ID               int       `gorm:"primaryKey" json:"id,omitempty"`
	UserID           int       `gorm:"column:user_id;not null" json:"user_id,omitempty"`
	EncryptedMessage string    `gorm:"column:encrypted_message;type:text;not null" json:"encrypted_message,omitempty"`
	Key              string    `gorm:"column:key;type:text;not null" json:"key,omitempty"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
}
