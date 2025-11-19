package repository

import "trinity/internal/models"

// Create — создание записи Кузнечик
func (db *PostgresDB) CreateKuznechik(kuz *models.Kuznechik) error {
	return db.Conn.Create(kuz).Error
}

// GetKuznechikByID — получение по ID
func (db *PostgresDB) GetKuznechikByID(id int) (*models.Kuznechik, error) {
	var kuz models.Kuznechik
	err := db.Conn.First(&kuz, id).Error
	if err != nil {
		return nil, err
	}
	return &kuz, nil
}

// GetKuznechikByUserID — все записи Кузнечик для пользователя
func (db *PostgresDB) GetKuznechikListByUserID(userID uint) ([]models.Kuznechik, error) {
	var kuzList []models.Kuznechik
	err := db.Conn.Where("user_id = ?", userID).Find(&kuzList).Error
	return kuzList, err
}

func (db *PostgresDB) GetKuznechikListByLogin(login string) ([]models.Kuznechik, error) {
	var user *models.User
	err := db.Conn.Preload("Kuznechiks").Omit("password").Where("login = ?", login).First(&user).Error
	return user.Kuznechiks, err
}

// DeleteKuznechik — удаление по ID
func (db *PostgresDB) DeleteKuznechik(id int) error {
	return db.Conn.Delete(&models.Kuznechik{}, id).Error
}
