package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	neturl "net/url"

	userapi "github.com/sh-miyoshi/hekate/pkg/apihandler/admin/v1/user"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
)

func (h *Handler) getUserID(projectName string, userName string) (string, error) {
	user, err := h.UserGetList(projectName, userName)
	if err != nil {
		return "", err
	}
	if len(user) != 1 {
		if len(user) == 0 {
			return "", fmt.Errorf("No such user")
		}
		return "", fmt.Errorf("Unexpect the number of user %s, expect 1, but got %d", userName, len(user))
	}
	return user[0].ID, nil
}

// UserAdd ...
func (h *Handler) UserAdd(projectName string, req *userapi.UserCreateRequest) (*userapi.UserGetResponse, error) {
	url := fmt.Sprintf("%s/adminapi/v1/project/%s/user", h.serverAddr, projectName)
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
		var res userapi.UserGetResponse
		if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			return nil, err
		}

		return &res, nil
	}

	message := ""
	var res errors.HTTPResponse
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
		return nil, fmt.Errorf("User %s is already exists", req.Name)
	case 500:
		return nil, fmt.Errorf("Internal server error occuered. Message: %s", message)
	}
	return nil, fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}

// UserDelete ...
func (h *Handler) UserDelete(projectName string, userName string) error {
	userID, err := h.getUserID(projectName, userName)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/adminapi/v1/project/%s/user/%s", h.serverAddr, projectName, userID)
	httpRes, err := h.request("DELETE", url, nil)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusNoContent {
		return nil
	}

	message := ""
	var res errors.HTTPResponse
	if err := json.NewDecoder(httpRes.Body).Decode(&res); err == nil {
		message = res.Error
	} else {
		message = "No messages."
	}

	switch httpRes.StatusCode {
	case 403:
		return fmt.Errorf("Loggined user did not have permission. Please login with other user")
	case 404:
		return fmt.Errorf("User %s in project %s is not found", userName, projectName)
	case 500:
		return fmt.Errorf("Internal server error occuered. Message: %s", message)
	}
	return fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}

// UserGetList ...
func (h *Handler) UserGetList(projectName string, userName string) ([]*userapi.UserGetResponse, error) {
	url := fmt.Sprintf("%s/adminapi/v1/project/%s/user", h.serverAddr, projectName)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Add("Authorization", fmt.Sprintf("bearer %s", h.accessToken))

	if userName != "" {
		httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		values := neturl.Values{}
		values.Set("name", userName)
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
		var res []*userapi.UserGetResponse
		if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			return nil, err
		}

		return res, nil
	}

	message := ""
	var res errors.HTTPResponse
	if err := json.NewDecoder(httpRes.Body).Decode(&res); err == nil {
		message = res.Error
	} else {
		message = "No messages."
	}

	switch httpRes.StatusCode {
	case 403:
		return nil, fmt.Errorf("Loggined user did not have permission. Please login with other user")
	case 404:
		return nil, fmt.Errorf("User %s in project %s is not found", userName, projectName)
	case 500:
		return nil, fmt.Errorf("Internal server error occuered. Message: %s", message)
	}
	return nil, fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}

// UserRoleAdd ...
func (h *Handler) UserRoleAdd(projectName string, userName string, roleName string, roleType model.RoleType) error {
	userID, err := h.getUserID(projectName, userName)
	if err != nil {
		return err
	}

	roleID, err := h.getRoleID(projectName, roleName, roleType)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/adminapi/v1/project/%s/user/%s/role/%s", h.serverAddr, projectName, userID, roleID)
	httpRes, err := h.request("POST", url, nil)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusNoContent {
		return nil
	}

	message := ""
	var res errors.HTTPResponse
	if err := json.NewDecoder(httpRes.Body).Decode(&res); err == nil {
		message = res.Error
	} else {
		message = "No messages."
	}

	switch httpRes.StatusCode {
	case 400:
		return fmt.Errorf("Invalid request. Message: %s", message)
	case 403:
		return fmt.Errorf("Loggined user did not have permission. Please login with other user")
	case 404:
		return fmt.Errorf("User %s in project %s is not found", userName, projectName)
	case 409:
		return fmt.Errorf("Role %s is already appended", roleName)
	case 500:
		return fmt.Errorf("Internal server error occuered. Message: %s", message)
	}
	return fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}

// UserRoleDelete ...
func (h *Handler) UserRoleDelete(projectName string, userName string, roleName string, roleType model.RoleType) error {
	userID, err := h.getUserID(projectName, userName)
	if err != nil {
		return err
	}

	roleID, err := h.getRoleID(projectName, roleName, roleType)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/adminapi/v1/project/%s/user/%s/role/%s", h.serverAddr, projectName, userID, roleID)
	httpRes, err := h.request("DELETE", url, nil)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusNoContent {
		return nil
	}

	message := ""
	var res errors.HTTPResponse
	if err := json.NewDecoder(httpRes.Body).Decode(&res); err == nil {
		message = res.Error
	} else {
		message = "No messages."
	}

	switch httpRes.StatusCode {
	case 400:
		return fmt.Errorf("Invalid request. Message: %s", message)
	case 403:
		return fmt.Errorf("Loggined user did not have permission. Please login with other user")
	case 404:
		return fmt.Errorf("User %s in project %s is not found", userName, projectName)
	case 500:
		return fmt.Errorf("Internal server error occuered. Message: %s", message)
	}
	return fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}

// UserChangePassword ...
func (h *Handler) UserChangePassword(projectName string, userName string, newPassword string) error {
	userID, err := h.getUserID(projectName, userName)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/adminapi/v1/project/%s/user/%s/reset-password", h.serverAddr, projectName, userID)
	body, _ := json.Marshal(&userapi.UserResetPasswordRequest{
		Password: newPassword,
	})

	httpRes, err := h.request("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusOK {
		return nil
	}

	message := ""
	var res errors.HTTPResponse
	if err := json.NewDecoder(httpRes.Body).Decode(&res); err == nil {
		message = res.Error
	} else {
		message = "No messages."
	}

	switch httpRes.StatusCode {
	case 400:
		return fmt.Errorf("Invalid request. Message: %s", message)
	case 403:
		return fmt.Errorf("Loggined user did not have permission. Please login with other user")
	case 404:
		return fmt.Errorf("User %s in project %s is not found", userName, projectName)
	case 500:
		return fmt.Errorf("Internal server error occuered. Message: %s", message)
	}
	return fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}

// UserUnlock ...
func (h *Handler) UserUnlock(projectName string, userName string) error {
	userID, err := h.getUserID(projectName, userName)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/adminapi/v1/project/%s/user/%s/unlock", h.serverAddr, projectName, userID)
	httpRes, err := h.request("POST", url, nil)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusNoContent {
		return nil
	}

	message := ""
	var res errors.HTTPResponse
	if err := json.NewDecoder(httpRes.Body).Decode(&res); err == nil {
		message = res.Error
	} else {
		message = "No messages."
	}

	switch httpRes.StatusCode {
	case 403:
		return fmt.Errorf("Loggined user did not have permission. Please login with other user")
	case 404:
		return fmt.Errorf("User %s in project %s is not found", userName, projectName)
	case 500:
		return fmt.Errorf("Internal server error occuered. Message: %s", message)
	}
	return fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}
