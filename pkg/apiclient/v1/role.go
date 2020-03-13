package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	roleapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/customrole"
)

// RoleAdd ...
func (h *Handler) RoleAdd(projectName string, req *roleapi.CustomRoleCreateRequest) (*roleapi.CustomRoleGetResponse, error) {
	u := fmt.Sprintf("%s/api/v1/project/%s/role", h.serverAddr, projectName)
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequest("POST", u, bytes.NewReader(body))
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
		var res roleapi.CustomRoleGetResponse
		if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			return nil, err
		}

		return &res, nil
	}
	return nil, fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}

// RoleDelete ...
func (h *Handler) RoleDelete(projectName string, roleName string) error {
	role, err := h.RoleGetList(projectName, roleName)
	if err != nil {
		return err
	}
	if len(role) != 1 {
		if len(role) == 0 {
			return fmt.Errorf("No such role")
		}
		return fmt.Errorf("Unexpect the number of role %s, expect 1, but got %d", roleName, len(role))
	}

	roleID := role[0].ID

	u := fmt.Sprintf("%s/api/v1/project/%s/role/%s", h.serverAddr, projectName, roleID)
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

// RoleGetList ...
func (h *Handler) RoleGetList(projectName string, roleName string) ([]*roleapi.CustomRoleGetResponse, error) {
	u := fmt.Sprintf("%s/api/v1/project/%s/role", h.serverAddr, projectName)
	httpReq, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Authorization", fmt.Sprintf("bearer %s", h.accessToken))

	if roleName != "" {
		httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		values := url.Values{}
		values.Set("name", roleName)
		httpReq.URL.RawQuery = values.Encode()
	}

	httpRes, err := h.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	switch httpRes.StatusCode {
	case 200:
		var res []*roleapi.CustomRoleGetResponse
		if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			return nil, err
		}

		return res, nil
	}
	return nil, fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}
