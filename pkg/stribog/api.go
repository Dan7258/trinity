package stribog

import "encoding/hex"

type EncryptedData struct {
	EncryptedMessage string `json:"encrypted_message"`
}

func HashingText(text string) *EncryptedData {
	hash := new(chunk)
	streebogCore([]byte(text), hash)
	ed := new(EncryptedData)
	ed.EncryptedMessage = hex.EncodeToString(hash[:])
	return ed
}
