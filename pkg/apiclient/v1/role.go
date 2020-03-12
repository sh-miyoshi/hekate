package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	roleapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/customrole"
)

// RoleAdd ...
func (h *Handler) RoleAdd(projectName string, req *roleapi.CustomRoleCreateRequest) (*roleapi.CustomRoleGetResponse, error) {
	url := fmt.Sprintf("%s/api/v1/project/%s/role", h.serverAddr, projectName)
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
		var res roleapi.CustomRoleGetResponse
		if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			return nil, err
		}

		return &res, nil
	}
	return nil, fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}

// RoleDelete ...
func (h *Handler) RoleDelete(projectName string, roleID string) error {
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
func (h *Handler) RoleGetList(projectName string) ([]*roleapi.CustomRoleGetResponse, error) {
	url := fmt.Sprintf("%s/api/v1/project/%s/role", h.serverAddr, projectName)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Authorization", fmt.Sprintf("bearer %s", h.accessToken))
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

// RoleGet ...
func (h *Handler) RoleGet(projectName, roleName string) (*roleapi.CustomRoleGetResponse, error) {
	url := fmt.Sprintf("%s/api/v1/project/%s/role/%s", h.serverAddr, projectName, roleName)
	httpReq, err := http.NewRequest("GET", url, nil)
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
