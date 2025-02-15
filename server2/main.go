package main

import (
	"log"
	"net/http"
	"os"
)

func main() {

	handler := func (w http.ResponseWriter, r *http.Request) {
		l := log.New(os.Stdout, "[server2]: ", log.Ldate|log.Ltime)
		l.Printf("Received request from %s\n", r.RemoteAddr)
		w.Write([]byte("Hello from server2!"))
	}

	http.HandleFunc("/", handler)

	log.Println("Starting server2 on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}