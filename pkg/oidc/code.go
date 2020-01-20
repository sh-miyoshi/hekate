package oidc

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"html/template"
	"net/http"
	"time"
)

type sessionInfo struct {
	VerifyCode  string
	ExpiresIn   time.Time
	BaseRequest *AuthRequest
}

var (
	expiresTimeSec uint64
	userLoginHTML  string
	loginSessions  map[string]*sessionInfo // key: verifyCode
)

// InitAuthCodeConfig ...
func InitAuthCodeConfig(authCodeExpiresTimeSec uint64, authCodeUserLoginFile string) {
	expiresTimeSec = authCodeExpiresTimeSec
	userLoginHTML = authCodeUserLoginFile
	loginSessions = make(map[string]*sessionInfo)
}

// GenerateAuthCode ...
func GenerateAuthCode(clientID string, redirectURL string, userID string) (string, error) {
	code := &model.AuthCode{
		CodeID:      uuid.New().String(),
		ClientID:    clientID,
		RedirectURL: redirectURL,
		ExpiresIn:   time.Now().Add(time.Second * time.Duration(expiresTimeSec)),
		UserID:      userID,
	}

	err := db.GetInst().NewAuthCode(code)

	return code.CodeID, err
}

// ValidateAuthCode ...
func ValidateAuthCode(codeID string) (*model.AuthCode, int, string) {
	code, err := db.GetInst().GetAuthCode(codeID)
	if err != nil {
		if err == model.ErrNoSuchCode {
			// TODO(revoke all token in code.UserID) <- SHOULD
			logger.Info("No such code %s", codeID)
			return nil, http.StatusBadRequest, "No such code"
		}
		logger.Error("Failed to get auth code: %+v", err)
		return nil, http.StatusInternalServerError, "Internal Server Error"
	}
	logger.Debug("Code: %v", code)

	if time.Now().Unix() >= code.ExpiresIn.Unix() {
		logger.Info("code %s is expired at %v", codeID, code.ExpiresIn)
		return nil, http.StatusBadRequest, "Code expired"
	}

	return code, http.StatusOK, "ok"
}

// WriteUserLoginPage ...
func WriteUserLoginPage(code string, projectName string, w http.ResponseWriter) {
	tpl := template.Must(template.ParseFiles(userLoginHTML))
	url := "/api/v1/project/" + projectName + "/openid-connect/login?login_verify_code=" + code

	d := map[string]string{
		"URL": url,
	}

	tpl.Execute(w, d)
}

// RegisterUserLoginSession ...
func RegisterUserLoginSession(req *AuthRequest) string {
	code := uuid.New().String()
	loginSessions[code] = &sessionInfo{
		VerifyCode:  code,
		ExpiresIn:   time.Now().Add(time.Second * time.Duration(expiresTimeSec)),
		BaseRequest: req,
	}
	return code
}

// UserLoginVerify ...
func UserLoginVerify(code string) (*AuthRequest, error) {
	if s, ok := loginSessions[code]; ok {
		// Get is only allowed once
		delete(loginSessions, code)

		now := time.Now().Unix()
		if now > s.ExpiresIn.Unix() {
			return nil, errors.New("Session already expired")
		}
		return s.BaseRequest, nil
	}
	return nil, errors.New("No such session")
}
