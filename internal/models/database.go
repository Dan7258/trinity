package models

import (
	"gorm.io/gorm"
)

type Model interface {
	GetConn() *gorm.DB
	ConnectToDatabase() error
	CreateUser(user *User) error
	CreateRSA(rsa *RSA) error
	CreateKuznechik(kuz *Kuznechik) error
	CreateStribog(stribog *Stribog) error
	GetUserByID(id uint) (*User, error)
	GetUserByLogin(login string) (*User, error)
	GetUserByTelegram(telegram string) (*User, error)
	GetUserWithPasswordByLogin(login string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id uint) error
	ListUsers() ([]User, error)
	GetKuznechikListByUserID(userID uint) ([]Kuznechik, error)
	GetKuznechikListByLogin(login string) ([]Kuznechik, error)
	GetRSAListByUserID(userID uint) ([]RSA, error)
	GetRSAListByLogin(login string) ([]RSA, error)
	GetStribogListByUserID(userID uint) ([]Stribog, error)
	GetStribogListByLogin(login string) ([]Stribog, error)
}
