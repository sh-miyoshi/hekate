package model

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tt := []struct {
		projectName   string
		expectSuccess bool
	}{
		{"project-ok", true},
		{"project-ng-str-!", false},
		// TODO(add more test case)
	}

	for _, target := range tt {
		prjInfo := ProjectInfo{
			Name: target.projectName,
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
