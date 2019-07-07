package token

import (
	"testing"
	"time"
)

func TestParseHTTPHeaderToken(t *testing.T) {
	InitConfig(time.Minute*5, "jwt-server", "testtest")
	validToken, _ := Generate()

	// Test Cases
	tt := []struct {
		token      string
		expectPass bool
	}{
		{"Bearer " + validToken, true},
		{"bearer " + validToken, true},
		{"", false},
		{validToken, false},
		{"bbbb " + validToken, false},
	}

	for _, tc := range tt {
		_, err := ParseHTTPHeaderToken(tc.token)
		if tc.expectPass && err != nil {
			t.Errorf("handler should pass with token %s, but got error %v", tc.token, err)
		}
		if !tc.expectPass && err == nil {
			t.Errorf("handler should not pass with token %s, but error is nil", tc.token)
		}
	}
}
