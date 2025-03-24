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
	http.Handle("/signup/", http.StripPrefix("/signup/", http.FileServer(http.Dir("frontend/signup"))))
	http.Handle("/login/", http.StripPrefix("/login/", http.FileServer(http.Dir("frontend/login"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))

	http.HandleFunc("/api/key", authservice.AuthMiddleware(false, keyservice.Handler))
	http.HandleFunc("/api/message", authservice.AuthMiddleware(false, messageservice.Handler))
	http.HandleFunc("/api/auth/register", authservice.RegisterHandler)
	http.HandleFunc("/api/auth/login", authservice.LoginHandler)

	http.HandleFunc("/", authservice.AuthMiddleware(true,
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Request for: %s", r.URL.Path)
			http.FileServer(http.Dir("frontend/chat")).ServeHTTP(w, r)
			log.Printf("Response sent")
		},
	))
	log.Println("Serving on https://localhost:8443")

	err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}
