package main

import (
	"log"
	"net/http"
	"zawie.io/e2e/backend/keyservice"
)

// TODO: Read values from environment variables
var certPath = "cert.pem"
var keyPath = "key.pem"

func main() {
	// Serve static files inside each page directory
	http.Handle("/chat/", http.StripPrefix("/chat/", http.FileServer(http.Dir("frontend/chat"))))
	http.HandleFunc("/api/key", keyservice.Handler)

	log.Println("Serving on https://localhost:8443")

	err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}
