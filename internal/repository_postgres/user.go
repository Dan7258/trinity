package repository_postgres

import "trinity/internal/models"

// Create — создание пользователя
func (db *PostgresDB) CreateUser(user *models.User) error {
	return db.Conn.Create(user).Error
}

// GetUserByID — получение по ID
func (db *PostgresDB) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := db.Conn.Omit("password").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByLogin — получение по логину
func (db *PostgresDB) GetUserByLogin(login string) (*models.User, error) {
	var user models.User
	err := db.Conn.Omit("password").Where("login = ?", login).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserWithPasswordByLogin — отдельно для аутентификации
func (db *PostgresDB) GetUserWithPasswordByLogin(login string) (*models.User, error) {
	var user models.User
	err := db.Conn.Where("login = ?", login).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *PostgresDB) GetUserByTelegram(telegram string) (*models.User, error) {
	var user models.User
	err := db.Conn.Omit("password").Where("telegram = ?", telegram).First(&user).Error
	return &user, err
}

func (db *PostgresDB) UpdateUser(user *models.User) error {
	return db.Conn.Save(user).Error
}

func (db *PostgresDB) DeleteUser(id uint) error {
	return db.Conn.Delete(&models.User{}, id).Error
}

// ListUser — список всех пользователей
func (db *PostgresDB) ListUsers() ([]models.User, error) {
	var users []models.User
	err := db.Conn.Omit("password").Find(&users).Error
	return users, err
}
