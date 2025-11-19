package repository

import "trinity/internal/models"

// CreateStribog — создание записи Стрибог (хеш)
func (db *PostgresDB) CreateStribog(stribog *models.Stribog) error {
	return db.Conn.Create(stribog).Error
}

// GetStribogByID — получение по ID
func (db *PostgresDB) GetStribogByID(id int) (*models.Stribog, error) {
	var s models.Stribog
	err := db.Conn.First(&s, id).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// GetStribogByUserID — все хеши Стрибог для пользователя
func (db *PostgresDB) GetStribogListByUserID(userID uint) ([]models.Stribog, error) {
	var hashList []models.Stribog
	err := db.Conn.Where("user_id = ?", userID).Find(&hashList).Error
	return hashList, err
}

func (db *PostgresDB) GetStribogListByLogin(login string) ([]models.Stribog, error) {
	var user *models.User
	err := db.Conn.Preload("Stribogs").Omit("password").Where("login = ?", login).First(&user).Error
	return user.Stribogs, err
}

// DeleteStribog — удаление по ID
func (db *PostgresDB) DeleteStribog(id int) error {
	return db.Conn.Delete(&models.Stribog{}, id).Error
}
