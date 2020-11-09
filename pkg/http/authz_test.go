package http

import (
	"testing"
)

func TestGetTokenFromHeader(t *testing.T) {
	token := "testtokenstring"
	tt := []struct {
		tokenString   string
		expectSuccess bool
	}{
		{
			"bearer " + token,
			true,
		},
		{
			"Bearer " + token,
			true,
		},
		{
			"" + token,
			false,
		},
		{
			"bearerdummy " + token,
			false,
		},
	}

	for _, target := range tt {
		_, err := getTokenFromHeader(target.tokenString)
		if target.expectSuccess && err != nil {
			t.Errorf("Parse token %s returns wrong status. got %v, want nil", target.tokenString, err)
		}
		if !target.expectSuccess && err == nil {
			t.Errorf("Parse token %s returns wrong status. got nil, want error", target.tokenString)
		}
	}
}
