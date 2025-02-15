package main

import (
	"log"
	"net/http"
	"os"
)

func main() {

	handler := func (w http.ResponseWriter, r *http.Request) {
		l := log.New(os.Stdout, "[server1]: ", log.Ldate|log.Ltime)
		l.Printf("Received request from %s\n", r.RemoteAddr)
		w.Write([]byte("Hello from server1!"))
	}

	http.HandleFunc("/", handler)

	log.Println("Starting server1 on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}