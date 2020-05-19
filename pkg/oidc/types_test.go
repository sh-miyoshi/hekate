package oidc

import (
	"testing"
)

func TestValidateResponseType(t *testing.T) {
	supported := []string{
		"code",
		"id_token",
		"token",
		"code id_token",
		"code token",
		"id_token token",
		"code id_token token",
	}

	tt := []struct {
		types    []string
		expectOK bool
	}{
		{
			[]string{"code"},
			true,
		},
		{
			[]string{"token", "code"},
			true,
		},
		{
			[]string{"code", "token", "id_token"},
			true,
		},
		{
			[]string{"code", "invalid"},
			false,
		},
	}

	for _, tc := range tt {
		err := validateResponseType(tc.types, supported)
		if tc.expectOK && err != nil {
			t.Errorf("validateResponseType returns wrong response. got %v, want nil", err)
		}
		if !tc.expectOK && err == nil {
			t.Errorf("validateResponseType returns wrong response. got nil, but want not nil")
		}
	}
}
