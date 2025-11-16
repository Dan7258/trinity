package models

import "time"

type Stribog struct {
	ID        int       `gorm:"primaryKey" json:"id,omitempty"`
	UserID    int       `gorm:"column:user_id;not null" json:"user_id,omitempty"`
	Hash      string    `gorm:"column:encrypted_message;type:text;not null" json:"encrypted_message,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
}
