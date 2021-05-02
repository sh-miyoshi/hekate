package secret

import (
	"strings"

	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/stretchr/stew/slice"
)

var (
	ErrPasswordPolicyFailed = errors.New("Password do not much policy", "Password do not much policy")

	errPasswordTooShort = errors.Append(ErrPasswordPolicyFailed, "too short")
	errSameAsUserName   = errors.Append(ErrPasswordPolicyFailed, "same as user name")
	errBlackListed      = errors.Append(ErrPasswordPolicyFailed, "is in black list")
	errNotContainChar   = errors.Append(ErrPasswordPolicyFailed, "does not contain required character")
)

const (
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits       = "0123456789"
	specialChars = "!#$%&'()-=^~|@`[{]}:*;+,.<>/?_"
)

// CheckPassword ...
func CheckPassword(userName, password string, policy model.PasswordPolicy) *errors.Error {
	// MinimumLength
	if policy.MinimumLength > 0 {
		// If minimum length value is valid, check password length
		if uint(len(password)) < policy.MinimumLength {
			err := errPasswordTooShort.Copy()
			err.SetDescription("password must be at least %d characters", policy.MinimumLength)
			return err
		}
	}

	// NotUserName
	if policy.NotUserName && userName == password {
		err := errSameAsUserName.Copy()
		err.SetDescription("password must not be the same as user name")
		return err
	}

	// Black List
	if slice.Contains(policy.BlackList, password) {
		err := errBlackListed.Copy()
		err.SetDescription("password is in black list")
		return err
	}

	// UseCharacter
	if policy.UseCharacter != "" {
		err := errNotContainChar.Copy()

		switch policy.UseCharacter {
		case model.CharacterTypeLower:
			if !strings.ContainsAny(password, lowerChars) {
				err.SetDescription("password does not contain lowercase letters")
				return err
			}
		case model.CharacterTypeUpper:
			if !strings.ContainsAny(password, upperChars) {
				err.SetDescription("password does not contain uppercase letters")
				return err
			}
		case model.CharacterTypeBoth:
			if !strings.ContainsAny(password, lowerChars) || !strings.ContainsAny(password, upperChars) {
				err.SetDescription("password does not contain lowercase or uppercase letters")
				return err
			}
		case model.CharacterTypeEither:
			if !strings.ContainsAny(password, lowerChars) && !strings.ContainsAny(password, upperChars) {
				err.SetDescription("password does not contain alphabets")
				return err
			}
		}
	}

	// UseDigit
	if policy.UseDigit {
		if !strings.ContainsAny(password, digits) {
			err := errNotContainChar.Copy()
			err.SetDescription("password does not contain digits")
			return err
		}
	}

	// UseSpecialCharacter
	if policy.UseSpecialCharacter {
		if !strings.ContainsAny(password, specialChars) {
			err := errNotContainChar.Copy()
			err.SetDescription("password does not contain special characters")
			return err
		}
	}

	return nil
}
