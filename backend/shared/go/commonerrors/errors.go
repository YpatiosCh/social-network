package commonerrors

import (
	"context"
	"errors"
	"runtime"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error represents a custom error type that includes classification, cause, and context.
// It implements the error interface and supports error wrapping and classification.
type Error struct {
	code      error  // Classification: ErrNotFound, ErrInternal, etc. Enusured to never be nil
	input     string // The input given to the func returning or wraping: args, structs.
	stack     string // The stack starting from the most undeliyng error and three levels up.
	err       error  // Cause: wrapped original error.
	publicMsg string // A message that will be displayed to clients.
}

// Returns a string of the full stack of errors. For each error the string contains:
//   - Error.code: Classification: ErrNotFound, ErrInternal, etc. Enusured to never be nil
//   - Error.input: The input given to the func returning or wraping: args, structs
//   - Error.err: The wraped error down the chain.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	var builder strings.Builder

	if e.code != nil {
		builder.WriteString(e.code.Error())
	}

	if e.input != "" {
		builder.WriteString(": ")
		builder.WriteString(e.input)
	}

	if e.stack != "" {
		builder.WriteString(": ")
		builder.WriteString(e.stack)
	}

	if e.err != nil {
		builder.WriteString(": ")
		builder.WriteString(e.err.Error())
	}
	return builder.String()
}

// Stringer method for loggers
func (e *Error) String() string {
	return e.Error()
}

// Returns the original most underlying error by calling Unwrap until next err is nil.
func GetSource(err error) string {
	for {
		u := errors.Unwrap(err)
		if u == nil {
			return err.Error()
		}
		err = u
	}
}

// Method for errors.Is parsing. Returns `Error.code` match.
func (e *Error) Is(target error) bool {
	return e.code == target
}

// Method for error.As parsing. Returns the `MediaError.Err`.
func (e *Error) Unwrap() error {
	return e.err
}

// Creates a new Error with code
func New(code error, err error, msg ...string) *Error {
	if err == nil {
		return nil
	}

	e := &Error{
		code:  parseCode(code),
		err:   err,
		stack: getStack(3, 3),
		input: getMsg(msg...),
	}
	return e
}

// Wrap creates a MediaError that classifies and optionally wraps an existing error.
//
// Usage:
//   - kind: the classification of the error (e.g., ErrFailed, ErrNotFound). If nil, ErrUnknownClass is used.
//   - err: the underlying error to wrap; if nil, Wrap returns nil.
//   - msg: optional context message describing where or why the error occurred.
//
// Behavior:
//   - If `err` is already a MediaError and `kind` is nil, it preserves the original Kind and optionally adds a new message.
//   - Otherwise, it creates a new MediaError with the specified Kind, Err, and message.
//   - The resulting MediaError supports errors.Is (matches Kind) and errors.As (type assertion) and preserves the wrapped cause.
//   - If kind is nil and the err is not media error or lacks kind then kind is set to ErrUnknownClass.
//
// It is recommended to only use nil kind if the underlying error is of type Error and its kind is not nil.
func Wrap(code error, err error, msg ...string) *Error {
	if err == nil {
		return nil
	}

	var ce *Error
	if errors.As(err, &ce) {
		// Wrapping an existing custom error
		e := &Error{
			code:      ce.code,
			err:       err,
			publicMsg: ce.publicMsg, // retain public message by default
		}

		if code != nil {
			e.code = parseCode(code)
		}
		if e.code == nil {
			e.code = ErrUnknown
		}

		e.input = getMsg(msg...)

		return e
	}

	if code == nil {
		code = ErrUnknown
	}

	e := &Error{
		code: code,
		err:  err,
		// stack: getStack(3, 2),
	}

	e.input = getMsg(msg...)

	return e
}

// Add a Public Message to be displayed on APIs and other public endpoints
//
// Usage:
//
//	 return Wrap(ErrUnauthorized, err, "token expired").
//		WithPublic("Authentication required")
func (e *Error) WithPublic(msg string) *Error {
	e.publicMsg = msg
	return e
}

// Returns *Error e  with error code c. If c fails validation e's code becomes ErrUnknown.
func (e *Error) WithCode(c error) *Error {
	e.code = parseCode(c)
	return e
}

// Helper mapper from error to grpc code.
func ToGRPCCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}

	// TODO: Check this
	// Propagate gRPC status errors
	// if st, ok := status.FromError(err); ok {
	// 	return st.Code()
	// }

	// Handle context errors
	if errors.Is(err, context.DeadlineExceeded) {
		return codes.DeadlineExceeded
	}
	if errors.Is(err, context.Canceled) {
		return codes.Canceled
	}

	// Handle your domain error
	var e *Error
	if errors.As(err, &e) {
		if code, ok := errorToGRPC[e.code]; ok {
			return code
		}
	}

	// 4. Fallback
	return codes.Unknown
}

// Coverts a grpc error to commonerrors Error type.
// The status code is converted to commonerrors type and the status message is wraped inside it as a new error.
// Optionaly a msg string is included for additional context.
// Usefull for downstream error parsing.
func ParseGrpcErr(err error, msg ...string) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	code := st.Code() // codes.NotFound, codes.Internal, etc.
	message := st.Message()

	if domainErr, ok := grpcToError[code]; ok {
		return Wrap(domainErr, errors.New(message), getMsg(msg...))
	}
	return Wrap(ErrUnknown, err, getMsg(msg...))
}

// Converts a commonerrors type Error to grpc status error. Handles context errors first.
// If the error passed is neither context error or Error unknown is returned.
func GRPCStatus(err error) error {
	if err == nil {
		return nil
	}

	// Propagate gRPC status errors
	if st, ok := status.FromError(err); ok {
		return st.Err()
	}

	// Handle context errors
	if errors.Is(err, context.DeadlineExceeded) {
		return status.Errorf(codes.DeadlineExceeded, "deadline exceeded")
	}
	if errors.Is(err, context.Canceled) {
		return status.Errorf(codes.Canceled, "request canceled")
	}

	// Handle domain error
	var e *Error
	if errors.As(err, &e) {
		msg := e.publicMsg
		if msg == "" {
			msg = "missing error message"
		}

		if code, ok := errorToGRPC[e.code]; ok {
			return status.Errorf(code, "service error: %v", msg)
		}
	}
	return status.Errorf(codes.Unknown, "unknown error")
}

// Returns c error if c is not nil and is a defined error in commonerrors else returns ErrUnknown
func parseCode(c error) error {
	_, ok := errorToGRPC[c]
	if c == nil || !ok {
		c = ErrUnknown
	}
	return c
}

func getMsg(msg ...string) string {
	if len(msg) > 0 {
		return msg[0]
	}
	return ""
}

func getStack(depth int, skip int) string {
	var builder strings.Builder
	builder.Grow(150)
	pc := make([]uintptr, depth)
	n := runtime.Callers(skip, pc)
	if n == 0 {
		return "(no caller data)"
	}
	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)
	for {
		frame, more := frames.Next()
		name := frame.Func.Name()
		start := strings.LastIndex(name, "/")
		builder.WriteString("by ")
		builder.WriteString(name[start+1:])
		builder.WriteString(" at ")
		builder.WriteString(strconv.Itoa(frame.Line))
		builder.WriteString("\n")
		if !more {
			break
		}
	}

	return builder.String()
}
