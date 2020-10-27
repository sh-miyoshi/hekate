package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"net/http/httputil"

	keysapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/keys"
	projectapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/project"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
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
	dump, _ := httputil.DumpRequest(httpReq, true)
	print.Debug("Project add method request\n---\n %s\n---\n", dump)
	httpRes, err := h.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusOK {
		var res projectapi.ProjectGetResponse
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
	case 409:
		return nil, fmt.Errorf("Project %s is already exists", req.Name)
	case 500:
		return nil, fmt.Errorf("Internal server error occuered. Message: %s", message)
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
	dump, _ := httputil.DumpRequest(httpReq, false)
	print.Debug("Project delete method request\n---\n %s---\n", dump)
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
		return fmt.Errorf("Do not have permission. Message: %s", message)
	case 404:
		return fmt.Errorf("Project %s is not found", projectName)
	case 500:
		return fmt.Errorf("Internal server error occuered. Message: %s", message)
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
	dump, _ := httputil.DumpRequest(httpReq, false)
	print.Debug("Project get list method request\n---\n %s---\n", dump)
	httpRes, err := h.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusOK {
		var res []*projectapi.ProjectGetResponse
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
	case 500:
		return nil, fmt.Errorf("Internal server error occuered. Message: %s", message)
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
	dump, _ := httputil.DumpRequest(httpReq, false)
	print.Debug("Project get method request\n---\n %s---\n", dump)
	httpRes, err := h.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusOK {
		var res projectapi.ProjectGetResponse
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
	case 403:
		return nil, fmt.Errorf("Loggined user did not have permission. Please login with other user")
	case 404:
		return nil, fmt.Errorf("Project %s is not found", projectName)
	case 500:
		return nil, fmt.Errorf("Internal server error occuered. Message: %s", message)
	}
	return nil, fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}

// ProjectUpdate ...
func (h *Handler) ProjectUpdate(projectName string, req *projectapi.ProjectPutRequest) error {
	url := fmt.Sprintf("%s/api/v1/project/%s", h.serverAddr, projectName)
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequest("PUT", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", fmt.Sprintf("bearer %s", h.accessToken))
	dump, _ := httputil.DumpRequest(httpReq, true)
	print.Debug("Project update method request\n---\n %s\n---\n", dump)
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
	case 400:
		return fmt.Errorf("Invalid request. Message: %s", message)
	case 403:
		return fmt.Errorf("Loggined user did not have permission. Please login with other user")
	case 404:
		return fmt.Errorf("Project %s is not found", projectName)
	case 500:
		return fmt.Errorf("Internal server error occuered. Message: %s", message)
	}
	return fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}

// ProjectKeysGet ...
func (h *Handler) ProjectKeysGet(projectName string) (*keysapi.KeysGetResponse, error) {
	url := fmt.Sprintf("%s/api/v1/project/%s/keys", h.serverAddr, projectName)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Add("Authorization", fmt.Sprintf("bearer %s", h.accessToken))
	dump, _ := httputil.DumpRequest(httpReq, false)
	print.Debug("Project Get List method request\n---\n %s---\n", dump)
	httpRes, err := h.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusOK {
		var res keysapi.KeysGetResponse
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
	case 403:
		return nil, fmt.Errorf("Loggined user did not have permission. Please login with other user")
	case 404:
		return nil, fmt.Errorf("Project %s is not found", projectName)
	case 500:
		return nil, fmt.Errorf("Internal server error occuered. Message: %s", message)
	}
	return nil, fmt.Errorf("Unexpected http response got. Message: %s", httpRes.Status)
}
