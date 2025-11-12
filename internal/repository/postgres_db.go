package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type PostgresDB struct {
	Conn *gorm.DB
}

func (db *PostgresDB) GetConn() *gorm.DB {
	return db.Conn
}

func (db *PostgresDB) ConnectToDatabase() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_SSLMODE"))
	var err error
	db.Conn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return err
}
