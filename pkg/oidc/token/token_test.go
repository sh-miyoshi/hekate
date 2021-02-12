package token

import (
	"net/http"
	"testing"
)

func TestGetFullIssuer(t *testing.T) {
	tt := []struct {
		url    string
		expect string
	}{
		{
			"http://localhost:18443/authapi/v1/project/master/.well-known/openid-configuration",
			"http://localhost:18443/authapi/v1/project/master",
		},
		{
			"http://localhost:18443/invalid/path",
			"http://localhost:18443",
		},
	}

	for _, target := range tt {
		req, _ := http.NewRequest("GET", target.url, nil)

		res := GetFullIssuer(req)
		if res != target.expect {
			t.Errorf("GetIssuer returns wrong string. got %s, want %s", res, target.expect)
		}
	}
}

func TestGetExpectIssuer(t *testing.T) {
	tt := []struct {
		url    string
		expect string
	}{
		{
			"http://localhost:18443/authapi/v1/project/master/.well-known/openid-configuration",
			"http://localhost:18443",
		},
	}

	for _, target := range tt {
		req, _ := http.NewRequest("GET", target.url, nil)

		res := GetExpectIssuer(req)
		if res != target.expect {
			t.Errorf("GetIssuer returns wrong string. got %s, want %s", res, target.expect)
		}
	}
}
