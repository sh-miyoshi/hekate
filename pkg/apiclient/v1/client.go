package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	clientapi "github.com/sh-miyoshi/hekate/pkg/apihandler/admin/v1/client"
	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// ClientAdd ...
func (h *Handler) ClientAdd(projectName string, req *clientapi.ClientCreateRequest) (*clientapi.ClientGetResponse, error) {
	url := fmt.Sprintf("%s/adminapi/v1/project/%s/client", h.serverAddr, projectName)
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
		var res clientapi.ClientGetResponse
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
		return nil, fmt.Errorf("Client %s is already exists", req.ID)
	case 500:
		return nil, fmt.Errorf("Internal server error occuered. Message: %s", message)
	}
	return nil, fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}

// ClientDelete ...
func (h *Handler) ClientDelete(projectName string, clientID string) error {
	url := fmt.Sprintf("%s/adminapi/v1/project/%s/client/%s", h.serverAddr, projectName, clientID)
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
		return fmt.Errorf("Client %s in project %s is not found", clientID, projectName)
	case 500:
		return fmt.Errorf("Internal server error occuered. Message: %s", message)
	}
	return fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}

// ClientGetList ...
func (h *Handler) ClientGetList(projectName string) ([]*clientapi.ClientGetResponse, error) {
	url := fmt.Sprintf("%s/adminapi/v1/project/%s/client", h.serverAddr, projectName)
	httpRes, err := h.request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusOK {
		var res []*clientapi.ClientGetResponse
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
		return nil, fmt.Errorf("Project %s is not found", projectName)
	case 500:
		return nil, fmt.Errorf("Internal server error occuered. Message: %s", message)
	}
	return nil, fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}

// ClientGet ...
func (h *Handler) ClientGet(projectName, clientID string) (*clientapi.ClientGetResponse, error) {
	url := fmt.Sprintf("%s/adminapi/v1/project/%s/client/%s", h.serverAddr, projectName, clientID)
	httpRes, err := h.request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusOK {
		var res clientapi.ClientGetResponse
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
	case 403:
		return nil, fmt.Errorf("Loggined user did not have permission. Please login with other user")
	case 404:
		return nil, fmt.Errorf("Client %s in project %s is not found", clientID, projectName)
	case 500:
		return nil, fmt.Errorf("Internal server error occuered. Message: %s", message)
	}
	return nil, fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}

// ClientUpdate ...
func (h *Handler) ClientUpdate(projectName, clientID string, req *clientapi.ClientPutRequest) error {
	url := fmt.Sprintf("%s/adminapi/v1/project/%s/client/%s", h.serverAddr, projectName, clientID)
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	httpRes, err := h.request("PUT", url, bytes.NewReader(body))
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
		return fmt.Errorf("Client %s in project %s is not found", clientID, projectName)
	case 500:
		return fmt.Errorf("Internal server error occuered. Message: %s", message)
	}
	return fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}
