package utils

import (
	"math/rand"
	"time"
)

func GenerateID() string {
	// Generate a ID with number and char
	const length = 8
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range length {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
