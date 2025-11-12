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
func (db *PostgresDB) GetStribogByUserID(userID int) ([]models.Stribog, error) {
	var hashList []models.Stribog
	err := db.Conn.Where("user_id = ?", userID).Find(&hashList).Error
	return hashList, err
}

// DeleteStribog — удаление по ID
func (db *PostgresDB) DeleteStribog(id int) error {
	return db.Conn.Delete(&models.Stribog{}, id).Error
}
