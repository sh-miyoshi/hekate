package http

import (
	"testing"
	"github.com/sh-miyoshi/jwt-server/pkg/token"
	"net/http"
)

func TestValidateAPIRequest(t *testing.T) {
	// Init TOken Handler
	token.InitConfig("jwt-server", "33c8bea55c2cb2bd1e6f18bf5699d5d8e6c42810c94b5ec85f909a415e03281d")
	// Token will expire at 2119/11/15
	validToken:="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIwYzU1NTE4Ni00MTU2LTRlNDctOGU4Zi1hYjZkNjYyMTMxMjIiLCJleHAiOjQ3Mjk0NTA4MzMsImp0aSI6IjUwY2IyZjQ0LTBjYzAtNGY1ZS1iOTA3LTMwZmRjNDA3Mzk4YSIsImlhdCI6MTU3NTg1MDgzMywiaXNzIjoiand0LXNlcnZlciIsInN1YiI6IjBjNTU1MTg2LTQxNTYtNGU0Ny04ZThmLWFiNmQ2NjIxMzEyMiIsInJvbGVzIjpbInJlYWQtY2x1c3RlciIsIndyaXRlLWNsdXN0ZXIiLCJyZWFkLXByb2plY3QiLCJ3cml0ZS1wcm9qZWN0IiwicmVhZC1yb2xlIiwid3JpdGUtcm9sZSIsInJlYWQtdXNlciIsIndyaXRlLXVzZXIiXX0.pE0YyiCsnF4iUgddCxnRK89gZ4M1xDrxdx3kn7QpVgY"

	tt := []struct {
		headerKey string
		headerValue string
		expectSuccess bool
	}{
		{
			"Authorization",
			"bearer "+validToken,
			true,
		},
	}

	for _, target := range tt {
		header := http.Header{}
		header.Add(target.headerKey, target.headerValue)

		_, err := ValidateAPIRequest(header)
		if target.expectSuccess && err != nil {
			t.Errorf("ValidateAPIRequest returns wrong status. got %v, want nil", err)
		}
		if !target.expectSuccess && err == nil {
			t.Errorf("ValidateAPIRequest returns wrong status. got nil, want error")
		}
	}
}