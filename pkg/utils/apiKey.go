package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

func GenerateApiKey() string {
	byteArr := make([]byte, 32)

	_, err := rand.Read(byteArr)
	if err != nil {
		log.Panic("error reading")
	}

	encoded := base64.URLEncoding.EncodeToString(byteArr)
	return encoded
}
