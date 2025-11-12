package models

import (
	"gorm.io/gorm"
)

type Model interface {
	GetConn() *gorm.DB
	ConnectToDatabase() error
}
