package pwpol

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/stretchr/stew/slice"
)

var (
	// ErrPasswordPolicyFailed ...
	ErrPasswordPolicyFailed = errors.New("possword do not much policy")

	// ErrPasswordTooShort ...
	ErrPasswordTooShort = errors.Wrap(ErrPasswordPolicyFailed, "too short")
	// ErrSameAsUserName ...
	ErrSameAsUserName = errors.Wrap(ErrPasswordPolicyFailed, "same as user name")
	// ErrBlackListed ...
	ErrBlackListed = errors.Wrap(ErrPasswordPolicyFailed, "is in black list")
	// ErrNotContainChar ...
	ErrNotContainChar = errors.Wrap(ErrPasswordPolicyFailed, "do not contain required character")
)

const (
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits       = "0123456789"
	specialChars = "!#$%&'()-=^~|@`[{]}:*;+,.<>/?_"
)

// CheckPassword ...
func CheckPassword(userName, password string, policy model.PasswordPolicy) error {
	// MinimumLength
	if policy.MinimumLength > 0 {
		// If minimum length value is valid, check password length
		if uint(len(password)) < policy.MinimumLength {
			return ErrPasswordTooShort
		}
	}

	// NotUserName
	if policy.NotUserName && userName == password {
		return ErrSameAsUserName
	}

	// Black List
	if slice.Contains(policy.BlackList, password) {
		return ErrBlackListed
	}

	// UseCharacter
	if policy.UseCharacter != "" {
		switch policy.UseCharacter {
		case model.CharacterTypeLower:
			if !strings.ContainsAny(password, lowerChars) {
				return ErrNotContainChar
			}
		case model.CharacterTypeUpper:
			if !strings.ContainsAny(password, upperChars) {
				return ErrNotContainChar
			}
		case model.CharacterTypeBoth:
			if !strings.ContainsAny(password, lowerChars) || !strings.ContainsAny(password, upperChars) {
				return ErrNotContainChar
			}
		case model.CharacterTypeEither:
			if !strings.ContainsAny(password, lowerChars) && !strings.ContainsAny(password, upperChars) {
				return ErrNotContainChar
			}
		}
	}

	// UseDigit
	if policy.UseDigit {
		if !strings.ContainsAny(password, digits) {
			return ErrNotContainChar
		}
	}

	// UseSpecialCharacter
	if policy.UseSpecialCharacter {
		if !strings.ContainsAny(password, specialChars) {
			return ErrNotContainChar
		}
	}

	return nil
}
