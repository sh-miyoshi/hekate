package util

import (
	"crypto/sha512"
	"fmt"
)

// CreateHash ...
func CreateHash(data string) string {
	byteHash := sha512.Sum512([]byte(data))
	return fmt.Sprintf("%x", byteHash)
}
