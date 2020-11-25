package util

import (
	"crypto/rand"
	"math/big"
)

const (
	// CharTypeDigit ...
	CharTypeDigit uint = 1 << iota
	// CharTypeLower ...
	CharTypeLower
	// CharTypeUpper ...
	CharTypeUpper
)

// RandomString ...
func RandomString(n int, typ uint) string {
	letter := ""
	if typ&CharTypeDigit != 0 {
		letter += "0123456789"
	}
	if typ&CharTypeLower != 0 {
		letter += "abcdefghijklmnopqrstuvwxyz"
	}
	if typ&CharTypeUpper != 0 {
		letter += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	if letter == "" {
		return ""
	}

	b := make([]byte, n)
	for i := 0; i < n; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letter))))
		b[i] += letter[n.Int64()]
	}
	return string(b)
}
