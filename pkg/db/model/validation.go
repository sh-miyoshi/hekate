package model

import (
	"github.com/asaskevich/govalidator"
	"regexp"
)

// ValidateProjectName ...
func ValidateProjectName(name string) bool {
	prjNameRegExp := regexp.MustCompile(`^[a-z][a-z0-9\-\.\_]{3,63}$`)
	return prjNameRegExp.MatchString(name)
}

// ValidateTokenSigningAlgorithm ...
func ValidateTokenSigningAlgorithm(signAlg string) bool {
	validAlgs := []string{
		"RS256",
	}

	for _, alg := range validAlgs {
		if signAlg == alg {
			return true
		}
	}
	return false
}

// ValidateLifeSpan ...
func ValidateLifeSpan(span uint) bool {
	return span >= 1
}

// ValidateClientID ...
func ValidateClientID(clientID string) bool {
	clientIDRegExp := regexp.MustCompile(`^[a-z][a-z0-9\-\.\_]{3,63}$`)
	return clientIDRegExp.MatchString(clientID)
}

// ValidateClientSecret ...
func ValidateClientSecret(secret string, accessType string) bool {
	if accessType != "confidential" {
		return true
	}

	if !(8 <= len(secret) && len(secret) < 256) {
		return false
	}
	return true
}

// ValidateClientAccessType ...
func ValidateClientAccessType(typ string) bool {
	allowedTypes := []string{
		"public",
		"confidential",
	}
	for _, t := range allowedTypes {
		if t == typ {
			return true
		}
	}
	return false
}

// ValidateUserName ...
func ValidateUserName(name string) bool {
	if !(3 <= len(name) && len(name) < 64) {
		return false
	}
	return true
}

// ValidateUserID ...
func ValidateUserID(id string) bool {
	return govalidator.IsUUID(id)
}

// ValidateSessionID ...
func ValidateSessionID(id string) bool {
	return govalidator.IsUUID(id)
}

// ValidateCustomRoleName ...
func ValidateCustomRoleName(name string) bool {
	if !(3 <= len(name) && len(name) < 64) {
		return false
	}
	return true
}
