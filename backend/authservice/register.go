package authservice

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/google/uuid"
)

var validAliasRegex = regexp.MustCompile(`^[a-zA-Z0-9_]*$`)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Auth register handler called", r.Method, r.URL.Path)
	w.Header().Set("Content-Type", "text/plain")

	err, alias, password := extractAuthFromRequest(r)
	if err != nil {
		fmt.Fprint(w, "Invalid Authorization Header")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Printf("Recieved registration request for \"%s\"", alias)

	if !validAliasRegex.MatchString(alias) {
		log.Printf("Alias \"%s\" is malformed!\n", alias)
		fmt.Fprint(w, "Alias must be alphanumeric")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	exists, err := CheckIfAliasExists(alias)
	if err != nil {
		fmt.Fprint(w, "Something went wrong checking alias uniquness")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if exists {
		log.Printf("Alias \"%s\" already exists!\n", alias)
		fmt.Fprint(w, "Alias already exists")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	salt := createSalt()

	insertUser(id, alias, salt, HashWithSalt(password, salt))
	initiateSession(w, alias)

	w.WriteHeader(http.StatusOK)
}

func createSalt() string {
	b := make([]byte, 32) // Create a byte slice of desired length
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x\n", b[:8])
}
