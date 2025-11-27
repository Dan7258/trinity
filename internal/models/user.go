package models

type User struct {
	ID         uint        `gorm:"primaryKey" json:"id,omitempty"`
	Login      string      `gorm:"size:255;unique;not null" json:"login,omitempty"`
	Role       string      `gorm:"size:255;not null" json:"role,omitempty"`
	Telegram   string      `gorm:"size:255;unique; not null" json:"telegram,omitempty"`
	Password   string      `gorm:"size:255;" json:"password,omitempty"`
	RSAs       []RSA       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Kuznechiks []Kuznechik `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Stribogs   []Stribog   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
