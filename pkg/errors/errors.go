package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error ...
type Error struct {
	privateMsgs      []string
	publicMsg        string
	httpResponseCode int
}

// Name        string `json:"error"`
// 	Description string `json:"error_description"`
// 	Code        int    `json:"status_code"`

// Error ...
func (e *Error) Error() string {
	return e.publicMsg
}

// GetHTTPStatusCode ...
func (e *Error) GetHTTPStatusCode() int {
	return e.httpResponseCode
}

// New ...
func New(format string, a ...interface{}) *Error {
	// TODO(runtime.caller)

	res := &Error{}
	msg := fmt.Sprintf(format, a...)
	if msg != "" {
		res.privateMsgs = append(res.privateMsgs, msg)
	}

	return res
}

// Append ...
func Append(err *Error, format string, a ...interface{}) *Error {
	// TODO(runtime.caller)

	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(format, a...)
	if msg != "" {
		err.privateMsgs = append(err.privateMsgs, msg)
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
		if len(all.privateMsgs) == 0 || len(err.privateMsgs) == 0 {
			return false
		}
		if all.privateMsgs[0] != err.privateMsgs[0] {
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
		// TODO(logger)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// TODO(Print)
// TODO(HTTPResponse)
