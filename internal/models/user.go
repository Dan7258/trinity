package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id,omitempty"`
	Login    string `gorm:"size:255;unique;not null" json:"login,omitempty"`
	Password string `gorm:"size:255;" json:"password,omitempty"`
}
