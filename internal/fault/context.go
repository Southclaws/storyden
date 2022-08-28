package fault

import (
	"errors"
	"fmt"
	"runtime"

	pkg_errors "github.com/pkg/errors"
)

// ErrorInfo represents an unwrapped chain of errors in a serialisable format
// which you can use with your structured logging solution or server response
// mechanism. It contains any contextual key-value data as well as a simple call
// stack style list of "points of interest" where errors have been wrapped with
// additional call-site specific context.
type ErrorInfo struct {
	Message string         `json:"message"`
	Values  map[string]any `json:"values,omitempty"`
	Trace   []Location     `json:"trace,omitempty"`
}

// Location is a point of interest in the code. Each one of these represents a
// location where an error has been wrapped in some way with additional context.
type Location struct {
	Message  string `json:"message"`
	Location string `json:"location"`
}

// ErrorData provides something similar to a stack trace, but much more minimal,
// focused and simple.
//
// A stack trace contains every single line of code in the call stack (including
// Go runtime internals) as well as a whole load of information about the super
// low level details of the state of the process at the time of the error. This
// information may be useful in some contexts but in a lot of cases in real
// world applications where business logic and data flow are among the primary
// focuses, this highly granular style of assembly-level stack trace can create
// noise and make finding the source of what often turns out to be application
// level rather than binary level errors take longer than is necessary.
//
// Instead, ErrorData provides a set of signposts, decorated with the available
// contextual information in a simple list. The Fault library provides a way to
// annotate errors as you ascend the call stack and this unrolls all that info
// in the error chain out to a simple list.
//
func ErrorData(err error) ErrorInfo {
	message := err.Error()
	values := make(map[string]any)
	trace := make([]Location, 0)

	for err != nil {
		if f, ok := err.(interface {
			withLocation
			withValue
		}); ok {
			if m := f.ErrorData(); m != nil {
				for k, v := range m {
					values[k] = v
				}
			}

			trace = append(trace, Location{
				Message:  err.Error(),
				Location: f.Location(),
			})
		}

		if f, ok := err.(stackTracer); ok {
			trace = append(trace, Location{
				Message:  err.Error(),
				Location: fmt.Sprintf("%s", f.StackTrace()[0]),
			})
		}

		err = errors.Unwrap(err)
	}

	return ErrorInfo{
		Message: message,
		Values:  values,
		Trace:   trace,
	}
}

type withLocation interface {
	Location() string
}

type withValue interface {
	ErrorData() map[string]string
}

type stackTracer interface {
	StackTrace() pkg_errors.StackTrace
}

// getLocation is called from within error constructors, so its always 2 levels
// deep in relation to the desired source code information.
func getLocation() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
	}

	return fmt.Sprintf("%s:%d", file, line)
}
