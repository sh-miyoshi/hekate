package oidc

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"html/template"
	"net/http"
	"time"
)

type sessionInfo struct {
	VerifyCode string
	ExpiresIn  time.Time
	// todo(client info)
}

var (
	expiresTimeSec = uint64(10 * 60) // default: 10 minutes
	userLoginHTML  = ""
)

// InitAuthCodeConfig ...
func InitAuthCodeConfig(authCodeExpiresTimeSec uint64, authCodeUserLoginFile string) {
	expiresTimeSec = authCodeExpiresTimeSec
	userLoginHTML = authCodeUserLoginFile
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
func WriteUserLoginPage(w http.ResponseWriter) {
	tpl := template.Must(template.ParseFiles(userLoginHTML))

	d := map[string]string{
		"URL": "http://localhost:8080",
	}

	tpl.Execute(w, d)
}

// RegisterUserLoginSession ...
func RegisterUserLoginSession() {
	// TODO(register to session list)
}

// UserLoginVerify ...
func UserLoginVerify(code string) error {
	// TODO(check session list)
	return nil
}
