package model

import (
	"regexp"
)

func validateProjectName(name string) bool {
	prjNameRegExp := regexp.MustCompile(`^[a-z][a-z0-9\-]{2,31}$`)
	return !prjNameRegExp.MatchString(name)
}

func validateTokenSigningAlgorithm(signAlg string) bool {
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

func validateLifeSpan(span uint) bool {
	return span >= 1
}

func validateClientID(clientID string) bool {
	if !(2 <= len(clientID) && len(clientID) < 128) {
		return false
	}
	return true
}

func validateClientSecret(secret string) bool {
	if !(8 <= len(secret) && len(secret) < 256) {
		return false
	}
	return true
}

func validateClientAccessType(typ string) bool {
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

func validateUserName(name string) bool {
	if !(3 <= len(name) && len(name) < 64) {
		return false
	}
	return true
}