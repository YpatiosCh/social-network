package commonerrors

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Returns c error if c is not nil and is a defined error
// in commonerrors else returns ErrUnknown
func parseCode(c error) error {
	if c == nil {
		c = ErrUnknown
	}
	_, ok := classToGRPC[c]
	if !ok {
		c = ErrUnknown
	}
	return c
}

// namedValue represents a value explicitly labeled with a name.
// It is used to associate structured input with a meaningful identifier
// when building error context.
type namedValue struct {
	name  string
	value any
}

// Named creates a namedValue wrapper.
//
// When passed to getInput (and ultimately error constructors),
// the name is rendered alongside the formatted value as:
//
//	<name> = <formatted value>
//
// This allows callers to explicitly label important inputs
// rather than relying on positional formatting.
func Named(name string, value any) namedValue {
	return namedValue{name: name, value: value}
}

// getInput formats a variadic list of inputs into a single string.
//
// Behavior:
//   - Each argument is rendered on its own line.
//   - If the argument is a namedValue, it is rendered as:
//     "<name> = <formatted value>"
//   - Otherwise, the argument is rendered using FormatValue directly.
//   - The final trailing newline is trimmed.
//
// This function is typically used to capture contextual input
// when creating or wrapping errors.
func getInput(args ...any) string {
	var b strings.Builder

	for _, arg := range args {
		switch v := arg.(type) {
		case namedValue:
			b.WriteString(fmt.Sprintf("%s = %s\n", v.name, FormatValue(v.value)))
		default:
			b.WriteString(FormatValue(arg))
			b.WriteString("\n")
		}
	}

	return strings.TrimRight(b.String(), "\n")
}

// FormatValue converts an arbitrary Go value into a readable, deterministic
// string representation suitable for error context, debugging, or logging.
//
// It is designed to be:
//   - Safe: panics during reflection are recovered and rendered as "<unprintable>"
//   - Recursive: nested structs, slices, arrays, and maps are expanded
//   - Cycle-aware: pointer cycles are detected and rendered as "<cycle>"
//   - Stringer-aware: values implementing fmt.Stringer are rendered using String()
//
// FormatValue is the public entry point. It initializes the recursion depth
// and the cycle-detection map, then delegates to formatValueIndented.
func FormatValue(v any) string {
	return formatValueIndented(v, 0, make(map[uintptr]bool))
}

// formatValueIndented recursively formats a value with indentation.
//
// Parameters:
//   - v:     the value being formatted
//   - depth: current recursion depth, used to compute indentation
//   - seen:  a map of pointer addresses used for cycle detection
//
// Behavior overview:
//
//  1. Nil handling
//     - A nil interface or nil pointer renders as "nil".
//
//  2. Interface unwrapping
//     - Interfaces are repeatedly unwrapped until a concrete value is reached.
//     - This ensures formatting is based on the underlying value, not the interface.
//
//  3. Pointer handling
//     - Nil pointers render as "nil".
//     - Non-nil pointers are tracked by address to detect cycles.
//     - Cycles render as "<cycle>" to avoid infinite recursion.
//     - The pointer is dereferenced and formatting continues on the element.
//
//  4. Stringer support
//     - If the value implements fmt.Stringer, String() is used.
//     - If the value itself does not implement Stringer but its address does,
//     the pointer receiver String() method is used.
//
//  5. Composite types
//     - Structs: rendered as a block with field names and indented values.
//     * Unexported fields are shown as "<unexported>".
//     - Maps: rendered as key-value pairs, one per line.
//     - Slices/arrays: rendered as an indexed list, one element per line.
//
//  6. Fallback
//     - All other kinds fall back to fmt.Sprintf("%v").
//
// Panic safety:
//   - Any panic encountered during reflection is recovered and rendered
//     as "<unprintable>" to avoid crashing error construction.
func formatValueIndented(v any, depth int, seen map[uintptr]bool) (out string) {
	defer func() {
		if r := recover(); r != nil {
			out = "<unprintable>"
		}
	}()

	if v == nil {
		return "nil"
	}

	// PROTOBUF

	// Protobuf Timestamp (pointer)
	if ts, ok := v.(*timestamppb.Timestamp); ok {
		if ts == nil {
			return "nil"
		}
		if ts.IsValid() {
			return ts.AsTime().Format(time.RFC3339)
		}
		return "<invalid timestamp>"
	}
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// Unwrap interfaces
	for val.Kind() == reflect.Interface {
		if val.IsNil() {
			return "nil"
		}
		val = val.Elem()
		typ = val.Type()
	}

	// Handle pointers (with cycle detection)
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return "nil"
		}
		ptr := val.Pointer()
		if seen[ptr] {
			return "<cycle>"
		}
		seen[ptr] = true
		return formatValueIndented(val.Elem().Interface(), depth, seen)
	}

	indent := strings.Repeat("   ", depth)
	nextIndent := strings.Repeat("   ", depth+1)

	stringerType := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

	// Value implements Stringer
	if typ.Implements(stringerType) {
		return val.Interface().(fmt.Stringer).String()
	}

	// Pointer implements Stringer
	if val.CanAddr() {
		ptrVal := val.Addr()
		if ptrVal.Type().Implements(stringerType) {
			return ptrVal.Interface().(fmt.Stringer).String()
		}
	}

	switch val.Kind() {

	case reflect.Struct:
		var b strings.Builder
		name := typ.Name()
		if name == "" {
			name = "struct"
		}

		b.WriteString(name + " {\n")

		for i := 0; i < val.NumField(); i++ {
			fieldType := typ.Field(i)
			fieldVal := val.Field(i)

			b.WriteString(nextIndent + fieldType.Name + ": ")

			if fieldVal.CanInterface() {
				b.WriteString(formatValueIndented(
					fieldVal.Interface(),
					depth+1,
					seen,
				))
			} else {
				b.WriteString("<unexported>")
			}
			b.WriteString("\n")
		}

		b.WriteString(indent + "}")
		return b.String()

	case reflect.Map:
		var b strings.Builder
		b.WriteString("map {\n")

		for _, key := range val.MapKeys() {
			b.WriteString(nextIndent)
			b.WriteString(fmt.Sprintf(
				"%v: %s\n",
				key.Interface(),
				formatValueIndented(val.MapIndex(key).Interface(), depth+1, seen),
			))
		}

		b.WriteString(indent + "}")
		return b.String()

	case reflect.Slice, reflect.Array:
		var b strings.Builder
		b.WriteString("[ ")

		for i := 0; i < val.Len(); i++ {
			b.WriteString(formatValueIndented(
				val.Index(i).Interface(),
				depth+1,
				seen,
			))
			if i < val.Len()-1 {
				b.WriteString(", ")
			}
		}

		b.WriteString(indent + " ]")
		return b.String()

	default:
		return fmt.Sprintf("%v", v)
	}
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
	var count int
	for {
		count++
		frame, more := frames.Next()
		name := frame.Function
		start := strings.LastIndex(name, "/")
		builder.WriteString("level ")
		builder.WriteString(strconv.Itoa(count))
		builder.WriteString(" -> ")
		builder.WriteString(name[start+1:])
		builder.WriteString(" at l. ")
		builder.WriteString(strconv.Itoa(frame.Line))
		if !more {
			break
		}
		builder.WriteString("\n          ")
	}

	return builder.String()
}

// Helper mapper from error to grpc code.
func GetCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}

	// Propagate gRPC status errors
	if st, ok := status.FromError(err); ok {
		return st.Code()
	}

	// Handle context errors
	if errors.Is(err, context.DeadlineExceeded) {
		return codes.DeadlineExceeded
	}
	if errors.Is(err, context.Canceled) {
		return codes.Canceled
	}

	// Handle domain error
	var e *Error
	if errors.As(err, &e) {
		if code, ok := classToGRPC[e.class]; ok {
			return code
		}
	}

	// 4. Fallback
	return codes.Unknown
}
