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
		{userName: "admin", password: "Admin1!", policy: allPol, expect: ErrPasswordTooShort},       // length
		{userName: "Admin1234!", password: "Admin1234!", policy: allPol, expect: ErrSameAsUserName}, // not user name
		{
			userName: "admin",
			password: "Admin1234!",
			policy:   model.PasswordPolicy{BlackList: []string{"Admin1234!"}},
			expect:   ErrBlackListed,
		}, // black list
		{
			userName: "admin",
			password: "ADMIN1234!",
			policy:   model.PasswordPolicy{UseCharacter: model.CharacterTypeLower},
			expect:   ErrNotContainChar,
		}, // use chars(lower)
		{
			userName: "admin",
			password: "admin1234!",
			policy:   model.PasswordPolicy{UseCharacter: model.CharacterTypeUpper},
			expect:   ErrNotContainChar,
		}, // use chars(upper)
		{
			userName: "admin",
			password: "admin1234!",
			policy:   model.PasswordPolicy{UseCharacter: model.CharacterTypeBoth},
			expect:   ErrNotContainChar,
		}, // use chars(both)
		{
			userName: "admin",
			password: "ADMIN1234!",
			policy:   model.PasswordPolicy{UseCharacter: model.CharacterTypeBoth},
			expect:   ErrNotContainChar,
		}, // use chars(both)
		{
			userName: "admin",
			password: "1234567!",
			policy:   model.PasswordPolicy{UseCharacter: model.CharacterTypeEither},
			expect:   ErrNotContainChar,
		}, // use chars(either)
		{
			userName: "admin",
			password: "AdminPasswd!",
			policy:   model.PasswordPolicy{UseDigit: true},
			expect:   ErrNotContainChar,
		}, // use digits
		{
			userName: "admin",
			password: "Admin1234",
			policy:   model.PasswordPolicy{UseSpecialCharacter: true},
			expect:   ErrNotContainChar,
		}, // use special chars
	}

	for _, tc := range tt {
		err := CheckPassword(tc.userName, tc.password, tc.policy)
		if err != tc.expect {
			t.Errorf("Check %v returns wrong status. got %v, want %v", tc, err, tc.expect)
		}
	}
}
