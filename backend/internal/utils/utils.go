package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func WriteJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if code == http.StatusNoContent {
		return nil
	}
	// fmt.Println("sending this:", v)
	return json.NewEncoder(w).Encode(v)
}

func ErrorJSON(w http.ResponseWriter, code int, msg string) {
	err := WriteJSON(w, code, map[string]string{"error": msg})
	if err != nil {
		fmt.Printf("Failed to send error message: %s, code: %d, %s\n", msg, code, err)
	}
}

func B64urlEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func B64urlDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

type ctxKey string

const (
	ClaimsKey        ctxKey = "jwtClaims"
	ReqId            ctxKey = "X-Request-Id"
	ReqActionDetails ctxKey = "X-Action-Details"
	ReqTimestamp     ctxKey = "X-Timestamp"
)

func RequestWithValue[T any](r *http.Request, key ctxKey, val T) *http.Request {
	ctx := context.WithValue(r.Context(), key, val)
	return r.WithContext(ctx)
}

func GetValue[T any](r *http.Request, key ctxKey) (T, bool) {
	v := r.Context().Value(key)
	if v == nil {
		fmt.Println("v is nil")
		var zero T
		return zero, false
	}
	c, ok := v.(T)
	if !ok {
		panic(1) // this should never happen, which is why I'm putting a panic here so that this mistake is obvious
	}
	return c, ok
}

func CompareTimeStrings(current, prev string) (int, error) {
	currentTime, err := time.Parse(time.RFC3339, current)
	if err != nil {
		return 0, fmt.Errorf("failed to parse current time: %w", err)
	}
	prevTime, err := time.Parse(time.RFC3339, prev)
	if err != nil {
		return 0, fmt.Errorf("failed to parse previous time: %w", err)
	}
	return int(currentTime.Sub(prevTime).Milliseconds()), nil
}

func ParseReqSignature(r *http.Request, w http.ResponseWriter) (reqId string, timestamp string) {
	reqId, _ = GetValue[string](r, ReqId)
	if strings.TrimSpace(reqId) == "" {
		ErrorJSON(w, http.StatusBadRequest, "missing X-Request-Id header")
		return
	}

	timestamp, _ = GetValue[string](r, ReqTimestamp)
	if strings.TrimSpace(timestamp) == "" {
		ErrorJSON(w, http.StatusBadRequest, "missing X-Timestamp header")
		return
	}

	return reqId, timestamp
}
