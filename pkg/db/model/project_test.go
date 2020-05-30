package model

import (
	"testing"
)

func TestValidatePasswordPolicy(t *testing.T) {
	tt := []struct {
		charType      string
		expectSuccess bool
	}{
		{"", true},
		{"lower", true},
		{"upper", true},
		{"both", true},
		{"either", true},
		{"ng", false},
	}

	for _, tc := range tt {
		p := PasswordPolicy{
			UseCharacter: CharacterType(tc.charType),
		}
		err := p.validate()

		if tc.expectSuccess && err != nil {
			t.Errorf("Passowrd policy validate %v returns wrong status. got %v, want nil", tc, err)
		}
		if !tc.expectSuccess && err == nil {
			t.Errorf("Passowrd policy validate %v returns wrong status. got nil, want error", tc)
		}
	}
}

func TestValidate(t *testing.T) {
	tt := []struct {
		projectName          string
		tokenSigningAlg      string
		accessTokenLifeSpan  uint
		refreshTokenLifeSpan uint
		expectSuccess        bool
	}{
		{"project-ok._", "RS256", 1, 1, true},
		{"project-ng-str-!", "RS256", 1, 1, false},
		{"project-ok", "invalid", 1, 1, false},
		{"pr", "RS256", 1, 1, false},
		{"project-name-too-long0123456789012345678901234567890123456789012", "RS256", 1, 1, false},
		{"0prject", "RS256", 1, 1, false},
	}

	for _, target := range tt {
		prjInfo := ProjectInfo{
			Name: target.projectName,
			TokenConfig: &TokenConfig{
				AccessTokenLifeSpan:  target.accessTokenLifeSpan,
				RefreshTokenLifeSpan: target.refreshTokenLifeSpan,
				SigningAlgorithm:     target.tokenSigningAlg,
			},
		}
		err := prjInfo.Validate()

		if target.expectSuccess && err != nil {
			t.Errorf("Validate %v returns wrong status. got %v, want nil", target, err)
		}
		if !target.expectSuccess && err == nil {
			t.Errorf("Validate %v returns wrong status. got nil, want error", target)
		}
	}
}
