package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

func generateToken(lenght int) string {

	bytes := make([]byte, lenght)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate token: %v", err)

	}
	return base64.URLEncoding.EncodeToString(bytes)

}
