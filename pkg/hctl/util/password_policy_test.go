package util

import (
	"reflect"
	"testing"

	projectapi "github.com/sh-miyoshi/hekate/pkg/apihandler/admin/v1/project"
)

func TestParsePasswordPolicies(t *testing.T) {
	tt := []struct {
		policies []string
		expect   projectapi.PasswordPolicy
		expectOK bool
	}{
		{
			policies: []string{},
			expectOK: true,
		},
		{
			policies: []string{"minLen=8"},
			expect: projectapi.PasswordPolicy{
				MinimumLength: 8,
			},
			expectOK: true,
		},
		{
			policies: []string{"minLen=8=8"},
			expectOK: false,
		},
		{
			policies: []string{"minLen=test"},
			expectOK: false,
		},
	}

	for _, tc := range tt {
		res, err := ParsePolicies(tc.policies)
		if err != nil && tc.expectOK {
			t.Errorf("parsePolicies returns wrong response. input: %v, got %v, want nil", tc.policies, err)
		}
		if err == nil && !tc.expectOK {
			t.Errorf("parsePolicies returns wrong response. input: %v, got nil, want error", tc.policies)
		}
		if !reflect.DeepEqual(res, tc.expect) {
			t.Errorf("parsePolicies returns wrong response. input: %v, got %v, want %v", tc.policies, res, tc.expect)
		}
	}
}
