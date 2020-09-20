package secret

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"

	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// Keys ...
type Keys struct {
	Public  []byte
	Private []byte
}

// GetSignKey ...
func GetSignKey(alg string) (*Keys, *errors.Error) {
	switch alg {
	case "RS256":
		key, err := rsa.GenerateKey(rand.Reader, 2048) // fixed key length is ok?
		if err != nil {
			return nil, errors.New("RSA key generate failed", "Failed to generate RSA private key: %v", err)
		}
		privateKey := x509.MarshalPKCS1PrivateKey(key)
		publicKey := x509.MarshalPKCS1PublicKey(&key.PublicKey)
		return &Keys{Public: publicKey, Private: privateKey}, nil
	}

	return nil, errors.New("Invalid algorithm", "Algorithm %s is not defined", alg)
}
