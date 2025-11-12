package kuznechik

import (
	"encoding/hex"
	"fmt"
)

type EncryptedData struct {
	EncryptedMessage string `json:"encrypted_message"`
	Key              string `json:"key"`
}

// splitToBlocks делит данные на блоки по 16 байт, добавляя padding при необходимости
func splitToBlocks(data []byte) [][]byte {
	padding := BLOCK_SIZE - (len(data) % BLOCK_SIZE)
	if padding == 0 {
		padding = BLOCK_SIZE
	}

	// Добавляем padding по стандарту PKCS#7
	for i := 0; i < padding; i++ {
		data = append(data, byte(padding))
	}

	blocks := make([][]byte, 0, len(data)/BLOCK_SIZE)
	for i := 0; i < len(data); i += BLOCK_SIZE {
		blocks = append(blocks, data[i:i+BLOCK_SIZE])
	}
	return blocks
}

// joinBlocks объединяет все зашифрованные блоки в одну строку (hex-представление)
func joinBlocks(blocks [][]byte) string {
	result := make([]byte, 0, len(blocks)*BLOCK_SIZE*2)
	for _, b := range blocks {
		for _, v := range b {
			result = append(result, fmt.Sprintf("%02x", v)...)
		}
	}
	return string(result)
}

// EncryptText — готовая функция шифрования текста
func EncryptText(text string) *EncryptedData {
	data := []byte(text)
	blocks := splitToBlocks(data)
	encryptedBlocks := make([][]byte, 0, len(blocks))
	key := GenerateRandomKey()
	var roundKeys [10]chunk
	GenRoundKeys(key, &roundKeys)

	for _, b := range blocks {
		c := bytesToChunk((*[BLOCK_SIZE]byte)(b))
		var out chunk
		kuznechikEncrypt(&roundKeys, &c, &out)
		outBytes := chunkToBytes(&out)
		encryptedBlocks = append(encryptedBlocks, outBytes[:])
	}
	ed := EncryptedData{}
	ed.EncryptedMessage = joinBlocks(encryptedBlocks)
	ed.Key = KeyToString(key)
	return &ed
}

// splitEncrypted делит строку с hex-зашифрованными данными обратно на блоки
func splitEncrypted(cipherHex string) [][]byte {
	data := make([]byte, len(cipherHex)/2)
	for i := 0; i < len(cipherHex); i += 2 {
		fmt.Sscanf(cipherHex[i:i+2], "%02x", &data[i/2])
	}

	blocks := make([][]byte, 0, len(data)/BLOCK_SIZE)
	for i := 0; i < len(data); i += BLOCK_SIZE {
		blocks = append(blocks, data[i:i+BLOCK_SIZE])
	}
	return blocks
}

// removePadding удаляет PKCS#7 padding
func removePadding(data []byte) []byte {
	if len(data) == 0 {
		return data
	}
	padLen := int(data[len(data)-1])
	if padLen > BLOCK_SIZE {
		return data
	}
	return data[:len(data)-padLen]
}

// DecryptText — готовая функция расшифровки текста
func DecryptText(data EncryptedData) *string {
	key := StringToKey(data.Key)
	var roundKeys [10]chunk
	GenRoundKeys(key, &roundKeys)

	blocks := splitEncrypted(data.EncryptedMessage)
	decrypted := make([]byte, 0, len(blocks)*BLOCK_SIZE)

	for _, b := range blocks {
		c := bytesToChunk((*[BLOCK_SIZE]byte)(b))
		var out chunk
		kuznechikDecrypt(&roundKeys, &c, &out)
		outBytes := chunkToBytes(&out)
		decrypted = append(decrypted, outBytes[:]...)
	}

	decrypted = removePadding(decrypted)
	message := string(decrypted)
	return &message
}

// KeyToString преобразует []uint8 ключ в hex-строку
func KeyToString(key []uint8) string {
	return hex.EncodeToString(key)
}

// StringToKey преобразует hex-строку обратно в []uint8
func StringToKey(s string) []uint8 {
	key, err := hex.DecodeString(s)
	if err != nil {
		panic("ошибка парсинга ключа из строки: " + err.Error())
	}
	fmt.Println(key)
	return key
}
