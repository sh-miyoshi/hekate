package token

import (
	"net/http"
	"testing"
)

func TestGetExpectIssuer(t *testing.T) {
	InitConfig(false, "testsecret")

	tt := []struct {
		url    string
		expect string
	}{
		{
			"http://localhost:8080/api/v1/project/master/.well-known/openid-configuration",
			"http://localhost:8080",
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
