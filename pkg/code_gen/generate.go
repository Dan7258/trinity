package code_gen

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateRandomNumber() (string, error) {
	bigNum := big.NewInt(100000)
	randomNumber, err := rand.Int(rand.Reader, bigNum)
	if err != nil {
		fmt.Println("Ошибка при генерации случайного числа:", err)
		return "", err
	}
	return fmt.Sprintf("%.6d", randomNumber.Int64()), nil
}
