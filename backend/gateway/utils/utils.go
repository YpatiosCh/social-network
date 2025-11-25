package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

func GetEnv(key string) (v string, err error) {
	if v = os.Getenv(key); v != "" {
		return v, err
	}
	return "", fmt.Errorf("env variable for %s not found", v)
}
