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
	GetKuznechikListByUserID(userID uint) ([]Kuznechik, error)
	GetRSAListByUserID(userID uint) ([]RSA, error)
	GetStribogListByUserID(userID uint) ([]Stribog, error)
}
