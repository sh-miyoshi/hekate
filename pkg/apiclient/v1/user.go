package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	userapi "github.com/sh-miyoshi/jwt-server/pkg/apihandler/v1/user"
)

// UserAdd ...
func (h *Handler) UserAdd(projectName string, req *userapi.UserCreateRequest) (*userapi.UserGetResponse, error) {
	url := fmt.Sprintf("%s/api/v1/project/%s/user", h.serverAddr, projectName)
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", fmt.Sprintf("bearer %s", h.accessToken))
	httpRes, err := h.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	switch httpRes.StatusCode {
	case 200:
		var res userapi.UserGetResponse
		if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			return nil, err
		}

		return &res, nil
	}
	return nil, fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}

// UserDelete ...
func (h *Handler) UserDelete(projectName string, userName string) error {
	user, err := h.UserGetList(projectName, userName)
	if err != nil {
		return err
	}
	if len(user) != 1 {
		if len(user) == 0 {
			return fmt.Errorf("No such user")
		}
		return fmt.Errorf("Unexpect the number of user %s, expect 1, but got %d", userName, len(user))
	}

	userID := user[0].ID
	u := fmt.Sprintf("%s/api/v1/project/%s/user/%s", h.serverAddr, projectName, userID)
	httpReq, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}
	httpReq.Header.Add("Authorization", fmt.Sprintf("bearer %s", h.accessToken))

	httpRes, err := h.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	switch httpRes.StatusCode {
	case 204:
		return nil
	}
	return fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}

// UserGetList ...
func (h *Handler) UserGetList(projectName string, userName string) ([]*userapi.UserGetResponse, error) {
	u := fmt.Sprintf("%s/api/v1/project/%s/user", h.serverAddr, projectName)
	httpReq, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Add("Authorization", fmt.Sprintf("bearer %s", h.accessToken))

	if userName != "" {
		httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		values := url.Values{}
		values.Set("name", userName)
		httpReq.URL.RawQuery = values.Encode()
	}

	httpRes, err := h.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	switch httpRes.StatusCode {
	case 200:
		var res []*userapi.UserGetResponse
		if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			return nil, err
		}

		return res, nil
	case 404:
		return nil, fmt.Errorf("No such user in the project")
	}
	return nil, fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}
