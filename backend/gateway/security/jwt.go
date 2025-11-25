package security

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"social-network/gateway/utils"
	"strings"
	"time"
	// "platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/utils"
)

// ======== Minimal custom JWT (HS256) using ONLY Go stdlib ========
// This implementation creates and validates compact JWTs (header.payload.signature)
// with HMAC-SHA256. It keeps things small and readable for learning/testing.
//
// Security notes:
// - Use a long, random secret from env/secret manager in production.
// - Keep token lifetimes short and rotate secrets periodically.
// - Consider adding jti (token id) and a store if you need revocation.
// - Clock skew tolerance is applied for nbf/exp checks.

// Claims represents the token payload. Add fields you need.
type Claims struct {
	UserId int64 `json:"user_id"`       // user ID
	Exp    int64 `json:"exp"`           // expiration (unix seconds)
	Iat    int64 `json:"iat"`           // issued-at (unix seconds)
	Nbf    int64 `json:"nbf,omitempty"` // not-before (unix seconds)
	// Roles  []string `json:"roles,omitempty"`
	// You can embed custom fields as needed, e.g. Email, TenantID, etc.
}

type ctxKey string

// Holds the keys to values on request context.
// Use
const (
	ClaimsKey        ctxKey = "jwtClaims"
	ReqId            ctxKey = "X-Request-Id"
	ReqActionDetails ctxKey = "X-Action-Details"
	ReqTimestamp     ctxKey = "X-Timestamp"
)

var (
	// Replace with a strong secret (32+ random bytes) from env in real apps.
	secret = []byte(func() string {
		s, err := utils.GetEnv("jwt-key")
		if err != nil {
			fmt.Println(err)
			return ""
		}
		return s
	}())
	// Allow a small clock skew when validating nbf/exp.
	clockSkew = 30 * time.Second
)

// CreateToken builds a signed JWT string with HS256.
func CreateToken(claims Claims) (string, error) {
	head := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, err := json.Marshal(head)
	if err != nil {
		return "", err
	}
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	headerEnc := B64urlEncode(headerJSON)
	payloadEnc := B64urlEncode(payloadJSON)
	unsigned := headerEnc + "." + payloadEnc

	sig := signHS256(unsigned, secret)
	return unsigned + "." + sig, nil
}

func signHS256(unsigned string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(unsigned))
	return B64urlEncode(h.Sum(nil))
}

// ParseAndValidate verifies the signature and time-based claims.
func ParseAndValidate(token string) (Claims, error) {
	var zero Claims
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return zero, errors.New("invalid token format")
	}
	unsigned := parts[0] + "." + parts[1]
	expected := signHS256(unsigned, secret)
	if !hmac.Equal([]byte(parts[2]), []byte(expected)) {
		return zero, errors.New("invalid signature")
	}

	payload, err := B64urlDecode(parts[1])
	if err != nil {
		return zero, fmt.Errorf("payload base64: %w", err)
	}
	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return zero, fmt.Errorf("payload json: %w", err)
	}

	now := time.Now()
	if claims.Nbf != 0 {
		if now.Add(clockSkew).Before(time.Unix(claims.Nbf, 0)) {
			return zero, errors.New("token not yet valid")
		}
	}
	if claims.Exp != 0 {
		if now.After(time.Unix(claims.Exp, 0).Add(clockSkew)) {
			return zero, errors.New("token expired")
		}
	}
	return claims, nil
}

func B64urlEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func B64urlDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

// Adds value val to r context with key 'key'
func RequestWithValue[T any](r *http.Request, key ctxKey, val T) *http.Request {
	ctx := context.WithValue(r.Context(), key, val)
	return r.WithContext(ctx)
}

// Get value T from request context with key 'key'
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
