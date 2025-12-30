package commonerrors

import (
	"errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrOK                 = errors.New("OK")
	ErrCanceled           = errors.New("CANCELED")
	ErrUnknown            = errors.New("UNKNOWN")
	ErrInvalidArgument    = errors.New("INVALID_ARGUMENT")
	ErrDeadlineExceeded   = errors.New("DEADLINE_EXCEEDED")
	ErrNotFound           = errors.New("NOT_FOUND")
	ErrAlreadyExists      = errors.New("ALREADY_EXISTS")
	ErrPermissionDenied   = errors.New("PERMISSION_DENIED")
	ErrResourceExhausted  = errors.New("RESOURCE_EXHAUSTED")
	ErrFailedPrecondition = errors.New("FAILED_PRECONDITION")
	ErrAborted            = errors.New("ABORTED")
	ErrOutOfRange         = errors.New("OUT_OF_RANGE")
	ErrUnimplemented      = errors.New("UNIMPLEMENTED")
	ErrInternal           = errors.New("INTERNAL")
	ErrUnavailable        = errors.New("UNAVAILABLE")
	ErrDataLoss           = errors.New("DATA_LOSS")
	ErrUnauthenticated    = errors.New("UNAUTHENTICATED")
)

var errorToGRPC = map[error]codes.Code{
	ErrOK:                 codes.OK,
	ErrCanceled:           codes.Canceled,
	ErrUnknown:            codes.Unknown,
	ErrInvalidArgument:    codes.InvalidArgument,
	ErrDeadlineExceeded:   codes.DeadlineExceeded,
	ErrNotFound:           codes.NotFound,
	ErrAlreadyExists:      codes.AlreadyExists,
	ErrPermissionDenied:   codes.PermissionDenied,
	ErrResourceExhausted:  codes.ResourceExhausted,
	ErrFailedPrecondition: codes.FailedPrecondition,
	ErrAborted:            codes.Aborted,
	ErrOutOfRange:         codes.OutOfRange,
	ErrUnimplemented:      codes.Unimplemented,
	ErrInternal:           codes.Internal,
	ErrUnavailable:        codes.Unavailable,
	ErrDataLoss:           codes.DataLoss,
	ErrUnauthenticated:    codes.Unauthenticated,
}

func ToGRPCCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	code, ok := errorToGRPC[err.(*Error).Kind]
	if ok {
		return code
	}
	return codes.Unknown
}
