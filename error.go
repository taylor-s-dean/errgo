package errgo

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"google.golang.org/grpc/codes"
)

// Error is a custom error type that contains the error message, file, function, line, code, and stack trace.
type Error struct {
	Err        error
	Message    string
	File       string
	Function   string
	Line       int
	Code       codes.Code
	StackTrace string
}

// MarshalJSON implements the json.Marshaler interface.
func (e Error) MarshalJSON() ([]byte, error) {
	var errJSON json.RawMessage
	var err error

	if e.Err == nil {
		errJSON = []byte("null")
	} else if _, ok := e.Err.(*Error); ok {
		errJSON, err = json.Marshal(e.Err)
	} else {
		errJSON, err = json.Marshal(e.Err.Error())
	}
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]interface{}{
		"error":       errJSON,
		"message":     e.Message,
		"file":        e.File,
		"function":    e.Function,
		"line":        e.Line,
		"code":        fmt.Sprintf("%s (%d)", e.Code.String(), e.Code),
		"stack_trace": e.StackTrace,
	})
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("%s (%d): %s", e.Code.String(), e.Code, e.Message)
}

// Unwrap returns the underlying error for use with errors.Is and errors.As.
func (e *Error) Unwrap() error {
	return e.Err
}

func getStackInfo(stackDepth int) (string, string, int) {
	result_file := "?"
	result_func := "?()"
	result_line := 0

	if pc, file, line, ok := runtime.Caller(stackDepth + 1); ok {
		result_file = filepath.Base(file)
		result_line = line
		if fn := runtime.FuncForPC(pc); fn != nil {
			_, dotName, _ := strings.Cut(filepath.Base(fn.Name()), ".")
			result_func = strings.TrimLeft(dotName, ".") + "()"
		}
	}
	return result_file, result_func, result_line
}

// Unwrap returns the underlying error.
func Unwrap(err error) *Error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*Error); ok {
		return e
	}

	return nil
}

// GetRoot returns the root error.
func GetRoot(err error) *Error {
	if err == nil {
		return nil
	}

	var root *Error
	for e := Unwrap(err); e != nil; e = Unwrap(e.Err) {
		root = e
	}

	return root
}

// ContainsError returns true if the error contains the given error.
func ContainsError(left error, right error) bool {
	if left == nil || right == nil {
		return false
	}

	rightStr := right.Error()
	if left.Error() == rightStr {
		return true
	}

	for e := Unwrap(left); e != nil; e = Unwrap(e.Err) {
		if e.Error() == rightStr {
			return true
		}
		if e.Err != nil && e.Err.Error() == rightStr {
			return true
		}
	}

	return false
}

// ContainsCode returns true if the error contains the given code.
func ContainsCode(left error, right codes.Code) bool {
	if left == nil {
		return false
	}

	rightStr := right.String()
	if strings.Contains(left.Error(), rightStr) {
		return true
	}

	for e := Unwrap(left); e != nil; e = Unwrap(e.Err) {
		if e.Code == right {
			return true
		}
		if e.Err != nil && strings.Contains(e.Err.Error(), rightStr) {
			return true
		}
	}

	return false
}

// UnwrapAll returns all the errors in the error chain.
func UnwrapAll(err error) []*Error {
	var result []*Error
	if err == nil {
		return result
	}

	for e := Unwrap(err); e != nil; e = Unwrap(e.Err) {
		result = append(result, e)
	}

	return result
}

// JSON returns the error as a JSON string.
func (e *Error) JSON() string {
	jsonBytes, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

// JSON returns the error as a JSON string.
func JSON(e error) string {
	if e == nil {
		return ""
	}

	err := Unwrap(e)
	if err != nil {
		return err.JSON()
	}

	return e.Error()
}

// Wrap wraps the given error with the given message and code.
func Wrap(err error, msg string, code codes.Code) *Error {
	file, function, line := getStackInfo(1)

	// Only capture stack trace if one doesn't already exist in the error chain
	var stackTrace string
	if e, ok := err.(*Error); ok && e.StackTrace != "" {
		stackTrace = e.StackTrace
	} else {
		stackTrace = string(debug.Stack())
	}

	return &Error{
		Err:        err,
		Message:    msg,
		File:       file,
		Function:   function,
		Line:       line,
		Code:       code,
		StackTrace: stackTrace,
	}
}
