package logout

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Logout ...
func Logout(serverAddr, projectName, refreshToken string) error {
	u := fmt.Sprintf("%s/api/v1/project/%s/openid-connect/revoke", serverAddr, projectName)

	form := url.Values{}
	form.Add("token_type_hint", "refresh_token")
	form.Add("token", refreshToken)
	body := strings.NewReader(form.Encode())
	httpReq, err := http.NewRequest("POST", u, body)
	if err != nil {
		return fmt.Errorf("Failed to create logout request: %v", err)
	}
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	httpRes, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("Failed to request server: %v", err)
	}
	defer httpRes.Body.Close()

	if 500 <= httpRes.StatusCode && httpRes.StatusCode < 600 {
		return fmt.Errorf("Logout request failed with server error: %d, so please try again", httpRes.StatusCode)
	}
	return nil
}
