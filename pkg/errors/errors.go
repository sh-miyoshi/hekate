package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"github.com/sh-miyoshi/hekate/pkg/logger"
)

// Error ...
type Error struct {
	privateInfo      []info
	publicMsg        string
	description      string
	httpResponseCode int
}

type info struct {
	msg   string
	fname string
	line  int
}

// HTTPResponse ...
type HTTPResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	State            string `json:"state"`
}

// Error ...
func (e *Error) Error() string {
	return e.publicMsg
}

// Copy ...
func (e *Error) Copy() *Error {
	res := Error{
		publicMsg:        e.publicMsg,
		httpResponseCode: e.httpResponseCode,
		description:      e.description,
		privateInfo:      append([]info{}, e.privateInfo...),
	}

	return &res
}

// StatusCode ...
func (e *Error) StatusCode() int {
	return e.httpResponseCode
}

// SetDescription ...
func (e *Error) SetDescription(format string, a ...interface{}) {
	e.description = fmt.Sprintf(format, a...)
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

	resErr := err.Copy()
	msg := fmt.Sprintf(format, a...)
	if msg != "" {
		resErr.privateInfo = append(resErr.privateInfo, info{
			msg:   msg,
			fname: fname,
			line:  line,
		})
	}

	return resErr
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

		return true
	}

	return false
}

// WriteToHTTP ...
func WriteToHTTP(w http.ResponseWriter, err *Error, statusCode int, state string) {
	res := HTTPResponse{
		Error:            err.publicMsg,
		ErrorDescription: err.description,
		State:            state,
	}

	w.Header().Add("Content-Type", "application/json")

	c := statusCode
	if c == 0 {
		if err.httpResponseCode == 0 {
			logger.Error("Failed to get status code in error response: %v", res)
			res.Error = "Internal Server Error"
			c = http.StatusInternalServerError
		} else {
			c = err.httpResponseCode
		}
	}

	w.WriteHeader(err.httpResponseCode)

	logger.Debug("Return http error: code %d, body %v", err.httpResponseCode, res)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Error("Failed to encode response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// RedirectWithOAuthError ...
func RedirectWithOAuthError(w http.ResponseWriter, err *Error, method, redirectURL, state string) {
	if err == nil {
		return
	}

	values := url.Values{}
	values.Set("error", err.publicMsg)
	if err.description != "" {
		values.Set("error_description", err.description)
	}
	if state != "" {
		values.Set("state", state)
	}

	req, _ := http.NewRequest(method, redirectURL, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = values.Encode()

	logger.Debug("Return OAuth error to %s: %v", req.URL.String(), values)
	http.Redirect(w, req, req.URL.String(), http.StatusFound)
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
