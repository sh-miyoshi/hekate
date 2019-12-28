package token

import (
	"net/http"
	"testing"
)

func TestGetIssuer(t *testing.T) {
	// TODO(init config)
	InitConfig("testsecret")

	tt := []struct {
		url    string
		expect string
	}{
		{
			"http://localhost:8080/api/v1/project/master/.well-known/openid-configuration",
			"http://localhost:8080/api/v1/project/master",
		},
	}

	for _, target := range tt {
		req, _ := http.NewRequest("GET", target.url, nil)

		res := GetIssuer(req)
		if res != target.expect {
			t.Errorf("GetIssuer returns wrong string. got %s, want %s", res, target.expect)
		}
	}
}
