package numbersapi

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

type APIResponse struct {
	Data string `json:"data"`
}

func getRandomNumber() int {
	return rand.Intn(100)
}

func numbersapi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	responseData := APIResponse{Data: fmt.Sprintf("%d", getRandomNumber())}

	err := json.NewEncoder(w).Encode(responseData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
func RunCLI() {
	mux := http.NewServeMux()
	na := http.HandlerFunc(numbersapi)
	mux.Handle("GET /rnd", na)
	// mux.HandleFunc("/", numbersapi)
	fmt.Println("Starting server on :80")
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		log.Fatal(err)
	}
}
