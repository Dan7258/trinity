package models

import (
	"gorm.io/gorm"
)

type Model interface {
	GetConn() *gorm.DB
	ConnectToDatabase() error
	CreateUser(user *User) error
	GetUserByID(id uint) (*User, error)
	GetUserByLogin(login string) (*User, error)
	GetUserWithPasswordByLogin(login string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id uint) error
	ListUsers() ([]User, error)
}
