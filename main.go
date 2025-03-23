package main

import (
	"log"
	"net/http"

	"zawie.io/e2e/backend/authservice"
	"zawie.io/e2e/backend/keyservice"
	"zawie.io/e2e/backend/messageservice"
)

// TODO: Read values from environment variables
var certPath = "cert.pem"
var keyPath = "key.pem"

func main() {
	// Serve static files inside each page directory
	http.Handle("/chat/", http.StripPrefix("/chat/", http.FileServer(http.Dir("frontend/chat"))))
	http.Handle("/signup/", http.StripPrefix("/signup/", http.FileServer(http.Dir("frontend/signup"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))

	http.HandleFunc("/api/key", keyservice.Handler)
	http.HandleFunc("/api/message", messageservice.Handler)
	http.HandleFunc("/api/auth/register", authservice.RegisterHandler)

	log.Println("Serving on https://localhost:8443")

	err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}
