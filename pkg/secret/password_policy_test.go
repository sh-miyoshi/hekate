package secret

import (
	"testing"

	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
)

func TestCheckPassword(t *testing.T) {
	allPol := model.PasswordPolicy{
		MinimumLength:       8,
		NotUserName:         true,
		BlackList:           []string{"invalid-password"},
		UseCharacter:        model.CharacterTypeBoth,
		UseDigit:            true,
		UseSpecialCharacter: true,
	}

	tt := []struct {
		userName string
		password string
		policy   model.PasswordPolicy
		expect   *errors.Error
	}{
		{expect: nil}, // if empty policy, return ok
		{userName: "admin", password: "Admin1234!", policy: allPol, expect: nil},                    // correct password
		{userName: "admin", password: "Admin1!", policy: allPol, expect: errPasswordTooShort},       // length
		{userName: "Admin1234!", password: "Admin1234!", policy: allPol, expect: errSameAsUserName}, // not user name
		{
			userName: "admin",
			password: "Admin1234!",
			policy:   model.PasswordPolicy{BlackList: []string{"Admin1234!"}},
			expect:   errBlackListed,
		}, // black list
		{
			userName: "admin",
			password: "ADMIN1234!",
			policy:   model.PasswordPolicy{UseCharacter: model.CharacterTypeLower},
			expect:   errNotContainChar,
		}, // use chars(lower)
		{
			userName: "admin",
			password: "admin1234!",
			policy:   model.PasswordPolicy{UseCharacter: model.CharacterTypeUpper},
			expect:   errNotContainChar,
		}, // use chars(upper)
		{
			userName: "admin",
			password: "admin1234!",
			policy:   model.PasswordPolicy{UseCharacter: model.CharacterTypeBoth},
			expect:   errNotContainChar,
		}, // use chars(both)
		{
			userName: "admin",
			password: "ADMIN1234!",
			policy:   model.PasswordPolicy{UseCharacter: model.CharacterTypeBoth},
			expect:   errNotContainChar,
		}, // use chars(both)
		{
			userName: "admin",
			password: "1234567!",
			policy:   model.PasswordPolicy{UseCharacter: model.CharacterTypeEither},
			expect:   errNotContainChar,
		}, // use chars(either)
		{
			userName: "admin",
			password: "AdminPasswd!",
			policy:   model.PasswordPolicy{UseDigit: true},
			expect:   errNotContainChar,
		}, // use digits
		{
			userName: "admin",
			password: "Admin1234",
			policy:   model.PasswordPolicy{UseSpecialCharacter: true},
			expect:   errNotContainChar,
		}, // use special chars
	}

	for _, tc := range tt {
		err := CheckPassword(tc.userName, tc.password, tc.policy)
		if !errors.Contains(err, tc.expect) {
			t.Errorf("Check %v returns wrong status. got %v, want %v", tc, err, tc.expect)
		}
	}
}
