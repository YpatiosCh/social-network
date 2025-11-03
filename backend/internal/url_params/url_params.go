package urlparams

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var (
	ErrBadValue         = errors.New("bad value passed to url parameter")
	ErrNotEnoughParams  = errors.New("missing needed url parameters")
	ErrUnknownValueType = errors.New("unknown value type passed")
	ErrEmptyNeededVal   = errors.New("neede parameter is missing its value")
)

type Params map[string]string // key = parameter name, value = target type

func ParseUrlParams(r *http.Request, need Params, optional Params) (map[string]any, error) {
	result := make(map[string]any)
	query := r.URL.Query()

	// Handle required params
	for key, targetType := range need {
		strVal := strings.TrimSpace(query.Get(key))
		if strVal == "" {
			return nil, fmt.Errorf("missing required parameter: %s", key)
		}
		anyVal, err := parseValueStrict(targetType, strVal)
		if err != nil {
			return nil, fmt.Errorf("invalid value for required parameter %s: %w", key, err)
		}
		// For slices, make sure we actually got at least one valid value
		if isSliceType(targetType) {
			if isEmptySlice(anyVal) {
				return nil, fmt.Errorf("required slice parameter %s is empty or malformed", key)
			}
		}
		result[key] = anyVal
	}

	// Handle optional params
	for key, targetType := range optional {
		strVal := strings.TrimSpace(query.Get(key))
		if strVal == "" {
			result[key] = defaultValue(targetType)
			continue
		}
		anyVal, err := parseValueLenient(targetType, strVal)
		if err != nil {
			fmt.Println("something wrong with parseValueLenient:", err.Error())
		}
		result[key] = anyVal
	}

	return result, nil
}

// --- helpers ---

func parseValueStrict(targetType string, stringValue string) (any, error) {
	return parseValueGeneric(targetType, stringValue, false)
}

func parseValueLenient(targetType string, stringValue string) (any, error) {
	return parseValueGeneric(targetType, stringValue, true)
}

func parseValueGeneric(targetType string, stringValue string, lenient bool) (any, error) {
	switch targetType {
	case "string":
		if stringValue == "" && !lenient {
			return nil, ErrEmptyNeededVal
		}
		return stringValue, nil

	case "int":
		v, err := strconv.Atoi(stringValue)
		if err != nil {
			if lenient {
				return 0, nil // fallback
			}
			return nil, ErrBadValue
		}
		return v, nil

	case "int64":
		v, err := strconv.ParseInt(stringValue, 10, 64)
		if err != nil {
			if lenient {
				return int64(0), nil
			}
			return nil, ErrBadValue
		}
		return v, nil
	case "bool":
		v, err := strconv.ParseBool(stringValue)
		if err != nil {
			if lenient {
				return false, nil // fallback to false for optional
			}
			return nil, ErrBadValue
		}
		return v, nil

	case "int slice":
		return parseIntSlice(stringValue, lenient)

	case "int64 slice":
		return parseInt64Slice(stringValue, lenient)

	default:
		return nil, ErrUnknownValueType
	}
}

func parseIntSlice(input string, lenient bool) ([]int, error) {
	parts := strings.Split(input, ",")
	result := make([]int, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		v, err := strconv.Atoi(part)
		if err != nil {
			if lenient {
				continue // skip bad value
			}
			return nil, ErrBadValue
		}
		result = append(result, v)
	}
	if len(result) == 0 && !lenient {
		return nil, ErrEmptyNeededVal
	}
	return result, nil
}

func parseInt64Slice(input string, lenient bool) ([]int64, error) {
	parts := strings.Split(input, ",")
	result := make([]int64, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		v, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			if lenient {
				continue
			}
			return nil, ErrBadValue
		}
		result = append(result, v)
	}
	if len(result) == 0 && !lenient {
		return nil, ErrEmptyNeededVal
	}
	return result, nil
}

func defaultValue(targetType string) any {
	switch targetType {
	case "string":
		return ""
	case "int":
		return 0
	case "int64":
		return int64(0)
	case "int slice":
		return []int{}
	case "int64 slice":
		return []int64{}
	default:
		return nil
	}
}

func isSliceType(t string) bool {
	return t == "int slice" || t == "int64 slice"
}

func isEmptySlice(val any) bool {
	switch v := val.(type) {
	case []int:
		return len(v) == 0
	case []int64:
		return len(v) == 0
	default:
		return false
	}
}

// type Params map[string]string

// // Fill in maps as needed, key is the name of the parameter and value the target type you want that value to be coersed into
// //
// // supports target types: "int", "int64" "string", "int slice", "int64 slice"
// func ParseUrlParams(r *http.Request, need Params, optional Params) (map[string]any, error) {
// 	newMap := make(map[string]any)

// 	for key, targetType := range need {
// 		strVal := r.URL.Query().Get(key)
// 		strVal = strings.TrimSpace(strVal)
// 		if strVal == "" {
// 			return nil, ErrEmptyNeededVal
// 		}
// 		anyVal, err := parseValue(targetType, strVal)
// 		if err != nil {
// 			return nil, err
// 		}
// 		newMap[key] = anyVal
// 	}

// 	if len(need) > len(newMap) {
// 		return nil, ErrNotEnoughParams
// 	}

// 	for key, targetType := range optional {
// 		strVal := r.URL.Query().Get(key)
// 		strVal = strings.TrimSpace(strVal)
// 		if strVal == "" {
// 			continue
// 		}
// 		anyVal, err := parseValue(targetType, strVal)
// 		if err != nil {
// 			return nil, err
// 		}
// 		newMap[key] = anyVal
// 	}

// 	return newMap, nil
// }

// func parseValue(targetType string, stringValue string) (any, error) {
// 	var value any
// 	switch targetType {
// 	case "string":
// 		return stringValue, nil
// 	case "int":
// 		value, err := strconv.Atoi(stringValue)
// 		if err != nil {
// 			return nil, ErrBadValue
// 		}
// 		return value, nil
// 	case "int64":
// 		value, err := strconv.Atoi(stringValue)
// 		if err != nil {
// 			return nil, ErrBadValue
// 		}
// 		return int64(value), nil
// 	case "int slice":
// 		slice := []int{}

// 		parts := strings.Split(stringValue, ",")

// 		for _, strVal := range parts {
// 			val, err := parseValue("int", strVal)
// 			if err != nil {
// 				return value, err
// 			}
// 			x, ok := val.(int)
// if !ok {
// 	panic(1)
// }
// 			slice = append(slice, x)
// 		}
// 		return slice, nil
// 	case "int64 slice":
// 		slice := []int64{}

// 		parts := strings.Split(stringValue, ",")

// 		for _, strVal := range parts {
// 			val, err := parseValue("int64", strVal)
// 			if err != nil {
// 				return value, err
// 			}
// 			x, ok := val.(int64)
// if !ok {
// 			panic(1)
// 		}
// 			slice = append(slice, x)
// 		}
// 		return slice, nil
// 	case "bool":
// 		v, err := strconv.ParseBool(stringValue)
// 		if err != nil {
// 			return value, err
// 		}
// 		return v, nil

// 	default:
// 		return value, ErrUnknownValueType
// 	}
// }
