package util

import "math/rand"

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
		b[i] += letter[rand.Intn(len(letter))]
	}
	return string(b)
}
