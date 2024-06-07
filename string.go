package d

import (
	"crypto/rand"
	"encoding/base64"
)

type String struct{}

// Used to generate random strings
// Example: randomKey, err := GenerateRandomString(32)
// Return value example: tFredJ-Ii5Eh0hQAHaJXSSz8Ffd7S6xTY2s-ZMxOLCM=
func (s String) GenerateRandomString(length int) (string, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(key), nil
}
