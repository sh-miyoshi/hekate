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
