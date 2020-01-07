package model

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tt := []struct {
		projectName     string
		tokenSigningAlg string
		expectSuccess   bool
	}{
		{"project-ok", "RS256", true},
		{"project-ng-str-!", "RS256", false},
		{"project-ok", "invalid", false},
		// TODO(add more test case)
	}

	for _, target := range tt {
		prjInfo := ProjectInfo{
			Name: target.projectName,
			TokenConfig: &TokenConfig{
				SigningAlgorithm: target.tokenSigningAlg,
			},
		}
		err := prjInfo.Validate()

		if target.expectSuccess && err != nil {
			t.Errorf("Validate returns wrong status. got %v, want nil", err)
		}
		if !target.expectSuccess && err == nil {
			t.Errorf("Validate returns wrong status. got nil, want error")
		}
	}
}
