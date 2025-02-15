package main

import (
	"log"
	"net/http"
	"os"
)

func main() {

	handler := func (w http.ResponseWriter, r *http.Request) {
		l := log.New(os.Stdout, "[server3]: ", log.Ldate|log.Ltime)
		l.Printf("Received request from %s\n", r.RemoteAddr)
		w.Write([]byte("Hello from server3!"))
	}

	http.HandleFunc("/", handler)

	log.Println("Starting server3 on :8083")
	log.Fatal(http.ListenAndServe(":8083", nil))
}