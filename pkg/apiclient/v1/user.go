package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

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
	fmt.Printf("project: %s, user: %s\n", projectName, userName)
	//url := fmt.Sprintf("%s/api/v1/project/%s/user", h.serverAddr, projectName)
	// TODO(get userid by name, delete user by id)
	return nil
}
