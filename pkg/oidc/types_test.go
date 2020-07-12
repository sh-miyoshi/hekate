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
			t.Errorf("validateResponseType returns wrong response. input: %v, got %v, want nil", tc.types, err)
		}
		if !tc.expectOK && err == nil {
			t.Errorf("validateResponseType returns wrong response. input: %v, got nil, but want not nil", tc.types)
		}
	}
}

func TestValidateprompt(t *testing.T) {
	tt := []struct {
		prompts  []string
		expectOK bool
	}{
		{
			prompts:  []string{},
			expectOK: true,
		},
		{
			prompts:  []string{"login"},
			expectOK: true,
		},
		{
			prompts:  []string{"select_account"},
			expectOK: true,
		},
		{
			prompts:  []string{"consent"},
			expectOK: true,
		},
		{
			prompts:  []string{"select_account", "login", "consent"},
			expectOK: true,
		},
		{
			prompts:  []string{"none"},
			expectOK: true,
		},
		{
			prompts:  []string{"invalid"},
			expectOK: false,
		},
		{
			prompts:  []string{"none", "login"},
			expectOK: false,
		},
	}

	for _, tc := range tt {
		err := validatePrompt(tc.prompts)
		if tc.expectOK && err != nil {
			t.Errorf("validateprompt returns wrong response. input: %v, got %v, want nil", tc.prompts, err)
		}
		if !tc.expectOK && err == nil {
			t.Errorf("validateprompt returns wrong response. input: %v, got nil, but want not nil", tc.prompts)
		}
	}
}

func TestValidateResponseMode(t *testing.T) {
	tt := []struct {
		mode     string
		expectOK bool
	}{
		{
			mode:     "",
			expectOK: true,
		},
		{
			mode:     "query",
			expectOK: true,
		},
		{
			mode:     "fragment",
			expectOK: true,
		},
		{
			mode:     "invalid",
			expectOK: false,
		},
		{
			mode:     "queryfragment",
			expectOK: false,
		},
	}

	for _, tc := range tt {
		err := validateResponseMode(tc.mode)
		if tc.expectOK && err != nil {
			t.Errorf("validateResponseMode returns wrong response. input: %v, got %v, want nil", tc.mode, err)
		}
		if !tc.expectOK && err == nil {
			t.Errorf("validateResponseMode returns wrong response. input: %v, got nil, but want not nil", tc.mode)
		}
	}
}
