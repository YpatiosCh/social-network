package ct

import (
	"fmt"
)

// =======================
// RedisKey
// =======================

type RedisKey string

const (
	BasicUserInfo RedisKey = "basic_user_info"
	Image         RedisKey = "img"
)

// Construct generates the full redis key based on the RedisKey type and provided arguments.
// For BasicUserInfo, it expects 1 argument: id. Result: "basic_user_info:<id>"
// For Image, it expects 2 arguments: variant, id. Result: "img_<variant>:<id>"
// Returns an error if the arguments count is incorrect or if the arguments are invalid.
func (k RedisKey) Construct(args ...any) (string, error) {
	switch k {
	case BasicUserInfo:
		if len(args) != 1 {
			return "", fmt.Errorf("invalid number of arguments for BasicUserInfo key: expected 1, got %d", len(args))
		}
		id, ok := args[0].(Id)
		if !ok {
			return "", fmt.Errorf("invalid argument type for BasicUserInfo key: expected ct.Id, got %T", args[0])
		}
		if err := id.Validate(); err != nil {
			return "", fmt.Errorf("invalid id for BasicUserInfo key: %w", err)
		}
		return fmt.Sprintf("%s:%d", k, id), nil
	case Image:
		if len(args) != 2 {
			return "", fmt.Errorf("invalid number of arguments for Image key: expected 2, got %d", len(args))
		}
		variant, ok := args[0].(FileVariant)
		if !ok {
			return "", fmt.Errorf("invalid argument type for invariant (1st arg) for Image key: expected ct.FileVariant, got %T", args[0])
		}
		if err := variant.Validate(); err != nil {
			return "", fmt.Errorf("invalid variant for Image key: %w", err)
		}

		id, ok := args[1].(Id)
		if !ok {
			return "", fmt.Errorf("invalid argument type for id (2nd arg) for Image key: expected ct.Id, got %T", args[1])
		}
		if err := id.Validate(); err != nil {
			return "", fmt.Errorf("invalid id for Image key: %w", err)
		}
		return fmt.Sprintf("%s_%s:%d", k, variant, id), nil
	default:
		return "", fmt.Errorf("unknown RedisKey type: %s", k)
	}
}
