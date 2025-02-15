package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func main() {
	num := 1
	for {
		fmt.Println("Total number of reqs: ", num)
		sendRequest("localhost:8080")
		num += 1
		// time.Sleep(1 * time.Second)
	}
}

func sendRequest(addr string) bool {
	url, err := url.Parse("http://" + addr)
	
	if err != nil { //server down
		fmt.Println("Server down", addr)
		return false
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url.String())

	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Server down", addr)
		return false
	}
	defer resp.Body.Close()

	return true
}