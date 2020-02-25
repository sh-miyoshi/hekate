package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	projectapi "github.com/sh-miyoshi/jwt-server/pkg/apihandler/v1/project"
)

// ProjectAdd ...
func (h *Handler) ProjectAdd(req *projectapi.ProjectCreateRequest) (*projectapi.ProjectGetResponse, error) {
	url := fmt.Sprintf("%s/api/v1/project", h.serverAddr)
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
		var res projectapi.ProjectGetResponse
		if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			return nil, err
		}

		return &res, nil
	case 403:
		return nil, fmt.Errorf("Loggined user did not have write-cluster role. Please login with other user")
	case 409:
		return nil, fmt.Errorf("Project %s is already exists", req.Name)
	}
	return nil, fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}

// ProjectDelete ...
func (h *Handler) ProjectDelete(projectName string) error {
	url := fmt.Sprintf("%s/api/v1/project/%s", h.serverAddr, projectName)
	httpReq, err := http.NewRequest("DELETE", url, nil)
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

// ProjectGetList ...
func (h *Handler) ProjectGetList() ([]*projectapi.ProjectGetResponse, error) {
	url := fmt.Sprintf("%s/api/v1/project", h.serverAddr)
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
		var res []*projectapi.ProjectGetResponse
		if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			return nil, err
		}

		return res, nil
	}
	return nil, fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}

// ProjectGet ...
func (h *Handler) ProjectGet(projectName string) (*projectapi.ProjectGetResponse, error) {
	url := fmt.Sprintf("%s/api/v1/project/%s", h.serverAddr, projectName)
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
		var res projectapi.ProjectGetResponse
		if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			return nil, err
		}

		return &res, nil
	}
	return nil, fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}
