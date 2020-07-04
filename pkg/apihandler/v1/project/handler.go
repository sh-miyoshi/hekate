package projectapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	jwthttp "github.com/sh-miyoshi/hekate/pkg/http"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/role"
)

// AllProjectGetHandler ...
//   require role: read-cluster
func AllProjectGetHandler(w http.ResponseWriter, r *http.Request) {
	// Authorize API Request
	if err := jwthttp.Authorize(r, "", role.ResCluster, role.TypeRead); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	projects, err := db.GetInst().ProjectGetList()
	if err != nil {
		logger.Error("Failed to get project list: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res := []ProjectGetResponse{}
	for _, prj := range projects {
		grantTypes := []string{}
		for _, t := range prj.AllowGrantTypes {
			grantTypes = append(grantTypes, t.String())
		}

		res = append(res, ProjectGetResponse{
			Name:      prj.Name,
			CreatedAt: prj.CreatedAt,
			TokenConfig: TokenConfig{
				AccessTokenLifeSpan:  prj.TokenConfig.AccessTokenLifeSpan,
				RefreshTokenLifeSpan: prj.TokenConfig.RefreshTokenLifeSpan,
				SigningAlgorithm:     prj.TokenConfig.SigningAlgorithm,
			},
			PasswordPolicy: PasswordPolicy{
				MinimumLength:       prj.PasswordPolicy.MinimumLength,
				NotUserName:         prj.PasswordPolicy.NotUserName,
				BlackList:           prj.PasswordPolicy.BlackList,
				UseCharacter:        string(prj.PasswordPolicy.UseCharacter),
				UseDigit:            prj.PasswordPolicy.UseDigit,
				UseSpecialCharacter: prj.PasswordPolicy.UseSpecialCharacter,
			},
			AllowGrantTypes: grantTypes,
		})
	}
	logger.Debug("Project List: %v", res)

	jwthttp.ResponseWrite(w, "AllProjectGetHandler", &res)
}

// ProjectCreateHandler ...
//   require role: write-cluster
func ProjectCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Authorize API Request
	if err := jwthttp.Authorize(r, "", role.ResCluster, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse Request
	var request ProjectCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Info("Failed to decode project create request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Set Allow Grant Type List
	grantTypes := []model.GrantType{}
	for _, t := range request.AllowGrantTypes {
		v, err := model.GetGrantType(t)
		if err != nil {
			logger.Info("Failed to get grant type %s: %v", t, err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		grantTypes = append(grantTypes, v)
	}

	// Create Project Entry
	project := model.ProjectInfo{
		Name:         request.Name,
		CreatedAt:    time.Now(),
		PermitDelete: true,
		TokenConfig: &model.TokenConfig{
			AccessTokenLifeSpan:  request.TokenConfig.AccessTokenLifeSpan,
			RefreshTokenLifeSpan: request.TokenConfig.RefreshTokenLifeSpan,
			SigningAlgorithm:     request.TokenConfig.SigningAlgorithm,
		},
		PasswordPolicy: model.PasswordPolicy{
			MinimumLength:       request.PasswordPolicy.MinimumLength,
			NotUserName:         request.PasswordPolicy.NotUserName,
			BlackList:           request.PasswordPolicy.BlackList,
			UseCharacter:        model.CharacterType(request.PasswordPolicy.UseCharacter),
			UseDigit:            request.PasswordPolicy.UseDigit,
			UseSpecialCharacter: request.PasswordPolicy.UseSpecialCharacter,
		},
		AllowGrantTypes: grantTypes,
	}

	// Create New Project
	if err := db.GetInst().ProjectAdd(&project); err != nil {
		if errors.Contains(err, model.ErrProjectAlreadyExists) {
			logger.Info("Project %s is already exists", request.Name)
			http.Error(w, "Project Already Exists", http.StatusConflict)
		} else if errors.Contains(err, model.ErrProjectValidateFailed) {
			logger.Info("Invalid project entry is specified: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
		} else {
			logger.Error("Failed to create project: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Return Response
	res := ProjectGetResponse{
		Name:      project.Name,
		CreatedAt: project.CreatedAt,
		TokenConfig: TokenConfig{
			AccessTokenLifeSpan:  project.TokenConfig.AccessTokenLifeSpan,
			RefreshTokenLifeSpan: project.TokenConfig.RefreshTokenLifeSpan,
			SigningAlgorithm:     project.TokenConfig.SigningAlgorithm,
		},
		PasswordPolicy: PasswordPolicy{
			MinimumLength:       project.PasswordPolicy.MinimumLength,
			NotUserName:         project.PasswordPolicy.NotUserName,
			BlackList:           project.PasswordPolicy.BlackList,
			UseCharacter:        string(project.PasswordPolicy.UseCharacter),
			UseDigit:            project.PasswordPolicy.UseDigit,
			UseSpecialCharacter: project.PasswordPolicy.UseSpecialCharacter,
		},
		AllowGrantTypes: request.AllowGrantTypes,
	}

	jwthttp.ResponseWrite(w, "ProjectCreateHandler", &res)
}

// ProjectDeleteHandler ...
//   require role: write-cluster
func ProjectDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResCluster, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := db.GetInst().ProjectDelete(projectName); err != nil {
		if errors.Contains(err, model.ErrNoSuchProject) || errors.Contains(err, model.ErrProjectValidateFailed) {
			logger.Info("Project %s is not found: %v", projectName, err)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else if errors.Contains(err, model.ErrDeleteBlockedProject) {
			logger.Info("Failed to delete blocked project: %v", err)
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else {
			logger.Error("Failed to delete project: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Return 204 (No content) for success
	w.WriteHeader(http.StatusNoContent)
	logger.Info("ProjectDeleteHandler method successfully finished")
}

// ProjectGetHandler ...
//   require role: read-project
func ProjectGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeRead); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Get Project
	project, err := db.GetInst().ProjectGet(projectName)
	if err != nil {
		if errors.Contains(err, model.ErrNoSuchProject) {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to get project: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	grantTypes := []string{}
	for _, t := range project.AllowGrantTypes {
		grantTypes = append(grantTypes, t.String())
	}

	// Return Response
	res := ProjectGetResponse{
		Name:      project.Name,
		CreatedAt: project.CreatedAt,
		TokenConfig: TokenConfig{
			AccessTokenLifeSpan:  project.TokenConfig.AccessTokenLifeSpan,
			RefreshTokenLifeSpan: project.TokenConfig.RefreshTokenLifeSpan,
			SigningAlgorithm:     project.TokenConfig.SigningAlgorithm,
		},
		PasswordPolicy: PasswordPolicy{
			MinimumLength:       project.PasswordPolicy.MinimumLength,
			NotUserName:         project.PasswordPolicy.NotUserName,
			BlackList:           project.PasswordPolicy.BlackList,
			UseCharacter:        string(project.PasswordPolicy.UseCharacter),
			UseDigit:            project.PasswordPolicy.UseDigit,
			UseSpecialCharacter: project.PasswordPolicy.UseSpecialCharacter,
		},
		AllowGrantTypes: grantTypes,
	}

	jwthttp.ResponseWrite(w, "ProjectGetHandler", &res)
}

// ProjectUpdateHandler ...
//   require role: write-project
func ProjectUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse Request
	var request ProjectPutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Info("Failed to decode project update request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get Previous Project Info
	project, err := db.GetInst().ProjectGet(projectName)
	if err != nil {
		if errors.Contains(err, model.ErrNoSuchProject) || errors.Contains(err, model.ErrProjectValidateFailed) {
			logger.Info("Project %s is not found: %v", projectName, err)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to get project: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Update Parameters
	project.TokenConfig.AccessTokenLifeSpan = request.TokenConfig.AccessTokenLifeSpan
	project.TokenConfig.RefreshTokenLifeSpan = request.TokenConfig.RefreshTokenLifeSpan
	project.TokenConfig.SigningAlgorithm = request.TokenConfig.SigningAlgorithm
	project.PasswordPolicy.MinimumLength = request.PasswordPolicy.MinimumLength
	project.PasswordPolicy.NotUserName = request.PasswordPolicy.NotUserName
	project.PasswordPolicy.BlackList = request.PasswordPolicy.BlackList
	project.PasswordPolicy.UseCharacter = model.CharacterType(request.PasswordPolicy.UseCharacter)
	project.PasswordPolicy.UseDigit = request.PasswordPolicy.UseDigit
	project.PasswordPolicy.UseSpecialCharacter = request.PasswordPolicy.UseSpecialCharacter
	project.AllowGrantTypes = []model.GrantType{}
	for _, t := range request.AllowGrantTypes {
		v, err := model.GetGrantType(t)
		if err != nil {
			logger.Info("Failed to get grant type %s: %v", t, err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		project.AllowGrantTypes = append(project.AllowGrantTypes, v)
	}

	// Update DB
	if err := db.GetInst().ProjectUpdate(project); err != nil {
		if errors.Contains(err, model.ErrProjectValidateFailed) {
			logger.Error("Project info validation failed: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
		} else {
			logger.Error("Failed to update project: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
	logger.Info("ProjectUpdateHandler method successfully finished")
}
