package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	roleapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/customrole"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
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
	dump, _ := httputil.DumpRequest(httpReq, true)
	print.Debug("Role add method request\n---\n %s\n---\n", dump)
	httpRes, err := h.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusOK {
		var res roleapi.CustomRoleGetResponse
		if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			return nil, err
		}

		return &res, nil
	}

	message := ""
	var res errors.HTTPError
	if err := json.NewDecoder(httpRes.Body).Decode(&res); err == nil {
		message = res.Error
	} else {
		message = "No messages."
	}

	switch httpRes.StatusCode {
	case 400:
		return nil, fmt.Errorf("Invalid request. Message: %s", message)
	case 403:
		return nil, fmt.Errorf("Loggined user did not have permission. Please login with other user")
	case 404:
		return nil, fmt.Errorf("Project %s is not found", projectName)
	case 409:
		return nil, fmt.Errorf("Role %s is already exists", req.Name)
	case 500:
		return nil, fmt.Errorf("Internal server error occuered. Message: %s", message)
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
	dump, _ := httputil.DumpRequest(httpReq, false)
	print.Debug("Role delete method request\n---\n %s\n---\n", dump)

	httpRes, err := h.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusNoContent {
		return nil
	}

	message := ""
	var res errors.HTTPError
	if err := json.NewDecoder(httpRes.Body).Decode(&res); err == nil {
		message = res.Error
	} else {
		message = "No messages."
	}

	switch httpRes.StatusCode {
	case 403:
		return fmt.Errorf("Loggined user did not have permission. Please login with other user")
	case 404:
		return fmt.Errorf("Role %s in project %s is not found", roleID, projectName)
	case 500:
		return fmt.Errorf("Internal server error occuered. Message: %s", message)
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

	dump, _ := httputil.DumpRequest(httpReq, false)
	print.Debug("Role get list method request\n---\n %s\n---\n", dump)
	httpRes, err := h.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusOK {
		var res []*roleapi.CustomRoleGetResponse
		if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			return nil, err
		}

		return res, nil
	}

	message := ""
	var res errors.HTTPError
	if err := json.NewDecoder(httpRes.Body).Decode(&res); err == nil {
		message = res.Error
	} else {
		message = "No messages."
	}

	switch httpRes.StatusCode {
	case 403:
		return nil, fmt.Errorf("Loggined user did not have permission. Please login with other user")
	case 404:
		return nil, fmt.Errorf("Role %s in project %s is not found", roleName, projectName)
	case 500:
		return nil, fmt.Errorf("Internal server error occuered. Message: %s", message)
	}
	return nil, fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}
