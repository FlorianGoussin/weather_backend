package helper

import (
	"crypto/rand"
	"encoding/base64"
)

func generateRandomBytes(n int) ([]byte, error) {
	// Create a byte slice of length 'n' to store the random bytes
	b := make([]byte, n)

	// Fill the byte slice with random bytes
	_, err := rand.Read(b)
	if err != nil {
			return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	// Generate 's' random bytes
	b, err := generateRandomBytes(s)
	if err != nil {
			return "", err
	}

	// Encode the random bytes to base64 string
	return base64.URLEncoding.EncodeToString(b), nil
}