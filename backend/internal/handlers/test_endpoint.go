package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func SendText(w http.ResponseWriter, r *http.Request) {
	//CORS headers, they are needed so that the app at port 8080 can talk to the server at 8081
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.WriteHeader(http.StatusOK)
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error reading body: ", err.Error())
	}
	fmt.Println("received: ", string(bytes))
	json.NewEncoder(w).Encode("Hello! I'm backend server!")
}
