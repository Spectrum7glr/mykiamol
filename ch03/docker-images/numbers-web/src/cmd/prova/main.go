package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"numbersweb"
	"time"
)

func main() {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8091/rnd", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var jsonresp numbersweb.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&jsonresp)
	if err != nil {
		log.Fatal(err)
	}
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	fmt.Println(jsonresp.Data)

}
