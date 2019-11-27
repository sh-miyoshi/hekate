package util

import (
	"fmt"
	"crypto/sha512"
)

// CreateHash ...
func CreateHash(data string) string {
	byteHash := sha512.Sum512([]byte(data))
	return fmt.Sprintf("%x", byteHash)
}