package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/sh-miyoshi/hekate/pkg/logger"
)

type info struct {
	msg   string
	fname string
	line  int
}

// Error ...
type Error struct {
	privateInfo      []info
	publicMsg        string
	httpResponseCode int
}

// Error ...
func (e *Error) Error() string {
	return e.publicMsg
}

// GetHTTPStatusCode ...
func (e *Error) GetHTTPStatusCode() int {
	return e.httpResponseCode
}

// New ...
func New(publicMsg string, privateMsg string, a ...interface{}) *Error {
	_, fname, line, _ := runtime.Caller(1)

	msg := fmt.Sprintf(privateMsg, a...)
	res := &Error{
		publicMsg: publicMsg,
		privateInfo: []info{
			{
				msg:   msg,
				fname: fname,
				line:  line,
			},
		},
	}

	return res
}

// Append ...
func Append(err *Error, format string, a ...interface{}) *Error {
	_, fname, line, _ := runtime.Caller(1)

	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(format, a...)
	if msg != "" {
		err.privateInfo = append(err.privateInfo, info{
			msg:   msg,
			fname: fname,
			line:  line,
		})
	}

	return err
}

// UpdatePublicMsg ...
func UpdatePublicMsg(err *Error, format string, a ...interface{}) *Error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(format, a...)
	err.publicMsg = msg
	return err
}

// Contains ...
func Contains(all, err *Error) bool {
	if all == nil || err == nil {
		return false
	}

	if all.publicMsg == err.publicMsg {
		if len(all.privateInfo) == 0 || len(err.privateInfo) == 0 {
			return false
		}
		if all.privateInfo[0].msg != err.privateInfo[0].msg {
			return false
		}
	}

	return true
}

// WriteOAuthError ...
func WriteOAuthError(w http.ResponseWriter, err *Error, state string) {
	res := map[string]interface{}{
		"error": err.publicMsg,
		// TODO(error_description)
		"state": state,
	}

	w.Header().Add("Content-Type", "application/json")

	// TODO(code == 0 -> panic)
	w.WriteHeader(err.httpResponseCode)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Error("Failed to encode response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// RedirectWithOAuthError ...
func RedirectWithOAuthError(w http.ResponseWriter, method, redirectURL string, err *Error, state string) {
	info := map[string]interface{}{
		"error": err.publicMsg,
		// TODO(error_description)
		"state": state,
	}
	body, _ := json.Marshal(info)
	req, _ := http.NewRequest(method, redirectURL, bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	http.Redirect(w, req, redirectURL, http.StatusFound)
}

// Print ...
func Print(err *Error) {
	if err == nil {
		_, fname, line, _ := runtime.Caller(1)
		logger.ErrorCustom("%s:%d [ERROR] nil", fname, line)
		return
	}

	for i := len(err.privateInfo) - 1; i >= 0; i-- {
		msg := err.privateInfo[i].msg
		if i != len(err.privateInfo)-1 {
			msg = "|- " + msg
		}
		logger.ErrorCustom("%s:%d [ERROR] %s", err.privateInfo[i].fname, err.privateInfo[i].line, msg)
	}
}

// PrintAsInfo ...
func PrintAsInfo(err *Error) {
	_, fname, line, _ := runtime.Caller(1)

	if err == nil {
		logger.ErrorCustom("%s:%d [INFO] nil", fname, line)
		return
	}

	msg := ""
	for _, info := range err.privateInfo {
		msg = info.msg + ": " + msg
	}
	msg = strings.TrimSuffix(msg, ": ")
	logger.ErrorCustom("%s:%d [INFO] %s", fname, line, msg)
}
