package main

import (
	"simplecall"
	"time"
)

func main() {
	client := simplecall.InitHttpClient(time.Duration(2 * time.Second))
	req, err := client.BuildRequest("GET", "http://localhost:8088", nil, nil)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 1000; i++ {
		client.Client.Do(req)
	}

}
