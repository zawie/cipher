package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, HTTPS!")
}

// TODO: Read values from environment variables
var certPath = "cert.pem"
var keyPath = "key.pem"

func main() {
	http.HandleFunc("/", handler)

	// Use self-signed certificate (or a real one if available)
	log.Println("Starting HTTPS server on https://localhost:8443")
	err := http.ListenAndServeTLS(":8443", certPath, keyPath, nil)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}
