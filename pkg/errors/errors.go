package errors

import "fmt"

// Error ...
type Error struct {
	privateMsgs []string
	publicMsg   string
}

// Error ...
func (e *Error) Error() string {
	return e.publicMsg
}

// New ...
func New(privateMsg, publicMsg string) *Error {
	// TODO(runtime.caller)

	res := &Error{
		publicMsg: publicMsg,
	}

	if privateMsg != "" {
		res.privateMsgs = []string{privateMsg}
	}

	return res
}

// AppendPrivateMsg ...
func AppendPrivateMsg(err *Error, format string, a ...interface{}) *Error {
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

// TODO(Print)
// TODO(HTTPResponse)
