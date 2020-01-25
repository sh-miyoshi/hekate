package oidc

import (
	"crypto/x509"
	"github.com/dvsekhvalnov/jose2go/base64url"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/util"
)

// JWKInfo is a struct for JSON Web Key(JWK) format defined in https://tools.ietf.org/html/rfc7517
type JWKInfo struct {
	KeyType      string `json:"kty"`
	KeyID        string `json:"kid"`
	Algorithm    string `json:"alg"`
	PublicKeyUse string `json:"use"`
	N            string `json:"n,omitempty"` // Use in RSA
	E            string `json:"e,omitempty"` // Use in RSA
	X            string `json:"x,omitempty"` // Use in EC
	Y            string `json:"y,omitempty"` // Use in EC
}

// JWKSet ...
type JWKSet struct {
	Keys []JWKInfo `json:"keys"`
}

// GenerateJWKSet ...
func GenerateJWKSet(signAlg string, publicKey []byte) (*JWKSet, error) {
	jwk := JWKInfo{
		KeyID:        uuid.New().String(),
		Algorithm:    signAlg,
		PublicKeyUse: "sig",
	}

	switch signAlg {
	case "RS256":
		jwk.KeyType = "RSA"
		key, err := x509.ParsePKCS1PublicKey(publicKey)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to parse RSA public key")
		}
		e := util.Int2bytes(uint64(key.E))
		jwk.E = base64url.Encode(e)
		jwk.N = base64url.Encode(key.N.Bytes())
	default:
		return nil, errors.New("Now such signing algorithm")
	}

	res := &JWKSet{}
	res.Keys = append(res.Keys, jwk)

	return res, nil
}
