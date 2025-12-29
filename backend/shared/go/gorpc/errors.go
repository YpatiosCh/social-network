package gorpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ClassifiedError struct {
	Class       ErrorClass
	GRPCCode    codes.Code
	Retryable   bool
	Description string
}

type ErrorClass string

const (
	ErrorClassClient       ErrorClass = "CLIENT_ERROR"
	ErrorClassServer       ErrorClass = "SERVER_ERROR"
	ErrorClassRetryable    ErrorClass = "RETRYABLE_DEPENDENCY_ERROR"
	ErrorClassNonRetryable ErrorClass = "NON_RETRYABLE_DEPENDENCY_ERROR"
	ErrorClassTimeout      ErrorClass = "TIMEOUT"
	ErrorClassCanceled     ErrorClass = "CANCELED"
	ErrorClassUnavailable  ErrorClass = "UNAVAILABLE"
	ErrorClassTransport    ErrorClass = "TRANSPORT_ERROR"
	ErrorClassUnknown      ErrorClass = "UNKNOWN_ERROR"
)

func Classify(err error) ClassifiedError {
	if err == nil {
		return ClassifiedError{}
	}

	// ---- Context errors (caller-side) ----
	if errors.Is(err, context.DeadlineExceeded) {
		return ClassifiedError{
			Class:       ErrorClassTimeout,
			GRPCCode:    codes.DeadlineExceeded,
			Retryable:   true,
			Description: "request timed out",
		}
	}

	if errors.Is(err, context.Canceled) {
		return ClassifiedError{
			Class:       ErrorClassCanceled,
			GRPCCode:    codes.Canceled,
			Retryable:   false,
			Description: "request canceled",
		}
	}

	// ---- gRPC status errors ----
	st, ok := status.FromError(err)
	if !ok {
		// Non-gRPC error (network / transport)
		return ClassifiedError{
			Class:       ErrorClassTransport,
			GRPCCode:    codes.Unknown,
			Retryable:   true,
			Description: "transport or connection error",
		}
	}

	code := st.Code()

	switch code {

	// ---- Client errors (do NOT retry) ----
	case codes.InvalidArgument,
		codes.NotFound,
		codes.AlreadyExists,
		codes.PermissionDenied,
		codes.Unauthenticated,
		codes.FailedPrecondition,
		codes.OutOfRange:

		return ClassifiedError{
			Class:       ErrorClassClient,
			GRPCCode:    code,
			Retryable:   false,
			Description: "downstream client error",
		}

	// ---- Server errors (retryable) ----
	case codes.Internal,
		codes.DataLoss,
		codes.Unknown:

		return ClassifiedError{
			Class:       ErrorClassServer,
			GRPCCode:    code,
			Retryable:   true,
			Description: "downstream server error",
		}

	// ---- Availability / dependency failures ----
	case codes.Unavailable:
		return ClassifiedError{
			Class:       ErrorClassUnavailable,
			GRPCCode:    code,
			Retryable:   true,
			Description: "downstream service unavailable",
		}

	case codes.ResourceExhausted:
		return ClassifiedError{
			Class:       ErrorClassRetryable,
			GRPCCode:    code,
			Retryable:   true,
			Description: "downstream resource exhausted",
		}

	case codes.Aborted:
		return ClassifiedError{
			Class:       ErrorClassRetryable,
			GRPCCode:    code,
			Retryable:   true,
			Description: "request aborted, retry recommended",
		}

	case codes.DeadlineExceeded:
		return ClassifiedError{
			Class:       ErrorClassTimeout,
			GRPCCode:    code,
			Retryable:   true,
			Description: "downstream deadline exceeded",
		}

	// ---- Explicit non-retryable dependency errors ----
	case codes.Unimplemented:
		return ClassifiedError{
			Class:       ErrorClassNonRetryable,
			GRPCCode:    code,
			Retryable:   false,
			Description: "method not implemented",
		}

	default:
		return ClassifiedError{
			Class:       ErrorClassUnknown,
			GRPCCode:    code,
			Retryable:   false,
			Description: "unclassified grpc error",
		}
	}
}
