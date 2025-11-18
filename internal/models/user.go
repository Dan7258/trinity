package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id,omitempty"`
	Login    string `gorm:"size:255;unique;not null" json:"login,omitempty"`
	Role     string `gorm:"size:255;not null" json:"role,omitempty"`
	Password string `gorm:"size:255;" json:"password,omitempty"`
}
