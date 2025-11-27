package repository_postgres

import "trinity/internal/models"

// Create — создание записи RSA
func (db *PostgresDB) CreateRSA(rsa *models.RSA) error {
	return db.Conn.Create(rsa).Error
}

// GetByID — получение по ID
func (db *PostgresDB) GetByID(id int) (*models.RSA, error) {
	var rsa models.RSA
	err := db.Conn.First(&rsa, id).Error
	if err != nil {
		return nil, err
	}
	return &rsa, nil
}

// GetByUserID — все RSA-записи пользователя
func (db *PostgresDB) GetRSAListByUserID(userID uint) ([]models.RSA, error) {
	var rsaList []models.RSA
	err := db.Conn.Where("user_id = ?", userID).Find(&rsaList).Error
	return rsaList, err
}

func (db *PostgresDB) GetRSAListByLogin(login string) ([]models.RSA, error) {
	var user *models.User
	err := db.Conn.Preload("RSAs").Omit("password").Where("login = ?", login).First(&user).Error
	return user.RSAs, err
}

// Delete — удаление по ID
func (db *PostgresDB) Delete(id int) error {
	return db.Conn.Delete(&models.RSA{}, id).Error
}
