package authservice

import (
	"fmt"
	"log"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Login handler called", r.Method, r.URL.Path)
	w.Header().Set("Content-Type", "text/plain")

	err, alias, password := extractAuthFromRequest(r)
	if err != nil {
		fmt.Fprint(w, "Invalid Authorization Header")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	salt, err := getSalt(alias)
	if err != nil {
		log.Printf("No salt found for alias \"%s\"\n", alias)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	hashedAndSaltedPassword := HashWithSalt(password, salt)

	log.Printf("Recieved login request for \"%s\"", alias)

	valid, err := checkPassword(alias, hashedAndSaltedPassword)

	if !valid {
		// TODO: Limit attempts
		log.Printf("Invalid password for \"%s\"\n", alias)
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	initiateSession(w, alias)

	w.WriteHeader(http.StatusOK)
}
