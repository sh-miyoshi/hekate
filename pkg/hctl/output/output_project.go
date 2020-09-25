package output

import (
	"encoding/json"
	"fmt"

	projectapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/project"
)

// ProjectInfoFormat ...
type ProjectInfoFormat struct {
	project *projectapi.ProjectGetResponse
}

// ProjectsInfoFormat ...
type ProjectsInfoFormat struct {
	projects []*projectapi.ProjectGetResponse
}

// NewProjectInfoFormat ...
func NewProjectInfoFormat(project *projectapi.ProjectGetResponse) *ProjectInfoFormat {
	return &ProjectInfoFormat{
		project: project,
	}
}

// NewProjectsInfoFormat ...
func NewProjectsInfoFormat(projects []*projectapi.ProjectGetResponse) *ProjectsInfoFormat {
	return &ProjectsInfoFormat{
		projects: projects,
	}
}

// ToText ...
func (f *ProjectInfoFormat) ToText() (string, error) {
	res := fmt.Sprintf("Name:                    %s\n", f.project.Name)
	res += fmt.Sprintf("Created Time:            %s\n", f.project.CreatedAt)
	res += fmt.Sprintf("Access Token Life Span:  %d [sec]\n", f.project.TokenConfig.AccessTokenLifeSpan)
	res += fmt.Sprintf("Refresh Token Life Span: %d [sec]\n", f.project.TokenConfig.RefreshTokenLifeSpan)
	res += fmt.Sprintf("Token Signing Algorithm: %s\n", f.project.TokenConfig.SigningAlgorithm)
	res += fmt.Sprintf("Allow Grant Types:       %v\n", f.project.AllowGrantTypes)
	res += fmt.Sprintf("Password Policies:\n")
	res += fmt.Sprintf("  Minimum Length:        %d\n", f.project.PasswordPolicy.MinimumLength)
	res += fmt.Sprintf("  Not User Name:         %v\n", f.project.PasswordPolicy.NotUserName)
	res += fmt.Sprintf("  Use Character:         %v\n", f.project.PasswordPolicy.UseCharacter)
	res += fmt.Sprintf("  Use Digit:             %v\n", f.project.PasswordPolicy.UseDigit)
	res += fmt.Sprintf("  Use Special Character: %v\n", f.project.PasswordPolicy.UseSpecialCharacter)
	res += fmt.Sprintf("  Black List:            %v\n", f.project.PasswordPolicy.BlackList)
	res += fmt.Sprintf("User Lock Enabled:       %v\n", f.project.UserLock.Enabled)
	res += fmt.Sprintf("Max Login Failure:       %d\n", f.project.UserLock.MaxLoginFailure)
	res += fmt.Sprintf("Lock Duration:           %d [sec]\n", f.project.UserLock.LockDuration)
	res += fmt.Sprintf("Failure Reset Time:      %d [sec]\n", f.project.UserLock.FailureResetTime)

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

// ToText ...
func (f *ProjectsInfoFormat) ToText() (string, error) {
	res := ""
	for i, prj := range f.projects {
		format := NewProjectInfoFormat(prj)
		msg, err := format.ToText()
		if err != nil {
			return "", err
		}
		res += msg
		if i < len(f.projects)-1 {
			res += "\n---\n"
		}
	}
	return res, nil
}

// ToJSON ...
func (f *ProjectsInfoFormat) ToJSON() (string, error) {
	bytes, err := json.Marshal(f.projects)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
