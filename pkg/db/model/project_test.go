package model

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tt := []struct {
		projectName          string
		tokenSigningAlg      string
		accessTokenLifeSpan  uint
		refreshTokenLifeSpan uint
		expectSuccess        bool
	}{
		{"project-ok", "RS256", 1, 1, true},
		{"project-ng-str-!", "RS256", 1, 1, false},
		{"project-ok", "invalid", 1, 1, false},
		{"pr", "RS256", 1, 1, false},
		{"project-name-too-long012345678901", "RS256", 1, 1, false},
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
