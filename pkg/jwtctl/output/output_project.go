package output

import (
	"encoding/json"
	"fmt"
	projectapi "github.com/sh-miyoshi/jwt-server/pkg/projectapi/v1"
)

// ProjectInfoFormat ...
type ProjectInfoFormat struct {
	project *projectapi.ProjectGetResponse
}

// NewProjectInfoFormat ...
func NewProjectInfoFormat(project *projectapi.ProjectGetResponse) *ProjectInfoFormat {
	return &ProjectInfoFormat{
		project: project,
	}
}

// ToText ...
func (f *ProjectInfoFormat) ToText() (string, error) {
	res := ""
	res += fmt.Sprintf("Name:                    %s", f.project.Name)
	res += fmt.Sprintf("Created Time:            %s", f.project.CreatedAt.String())
	res += fmt.Sprintf("Access Token Life Span:  %d [sec]", f.project.TokenConfig.AccessTokenLifeSpan)
	res += fmt.Sprintf("Refresh Token Life Span: %d [sec]", f.project.TokenConfig.RefreshTokenLifeSpan)
	return res, nil
}

// ToJSON ...
func (f *ProjectInfoFormat) ToJSON() (string, error) {
	bytes, err := json.Marshal(f.project)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
