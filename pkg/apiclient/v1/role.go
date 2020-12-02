package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	neturl "net/url"

	roleapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/customrole"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
)

func (h *Handler) getRoleID(projectName, roleName string, roleType model.RoleType) (string, error) {
	if roleType == model.RoleSystem {
		return roleName, nil
	}

	role, err := h.RoleGetList(projectName, roleName)
	if err != nil {
		return "", err
	}
	if len(role) != 1 {
		if len(role) == 0 {
			return "", fmt.Errorf("No such role")
		}
		return "", fmt.Errorf("Unexpect the number of role %s, expect 1, but got %d", roleName, len(role))
	}

	return role[0].ID, nil
}

// RoleAdd ...
func (h *Handler) RoleAdd(projectName string, req *roleapi.CustomRoleCreateRequest) (*roleapi.CustomRoleGetResponse, error) {
	url := fmt.Sprintf("%s/api/v1/project/%s/role", h.serverAddr, projectName)
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpRes, err := h.request("POST", url, bytes.NewReader(body))
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
	roleID, err := h.getRoleID(projectName, roleName, model.RoleCustom)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/v1/project/%s/role/%s", h.serverAddr, projectName, roleID)
	httpRes, err := h.request("DELETE", url, nil)
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
	url := fmt.Sprintf("%s/api/v1/project/%s/role", h.serverAddr, projectName)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Authorization", fmt.Sprintf("bearer %s", h.accessToken))

	if roleName != "" {
		httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		values := neturl.Values{}
		values.Set("name", roleName)
		httpReq.URL.RawQuery = values.Encode()
	}

	dump, _ := httputil.DumpRequest(httpReq, false)
	print.Debug("server request dump: %q", dump)

	httpRes, err := h.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	dump, _ = httputil.DumpResponse(httpRes, false)
	print.Debug("server response dump: %q", dump)

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
