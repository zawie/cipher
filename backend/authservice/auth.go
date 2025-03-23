package authservice

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"crypto/sha256"
	"encoding/base64"

	"github.com/google/uuid"
)

const AUTH_COOKIE_NAME string = "auth_token"

func extractAuthFromRequest(r *http.Request) (err error, alias string, password string) {
	log.Println("Auth register handler called", r.Method, r.URL.Path)

	authHeader, ok := r.Header["Authorization"]
	if !ok || len(authHeader) != 1 {
		err = errors.New("Missing or too many Authorization Headers")
		return
	}

	// Remove "Basic " prefix
	encoded := strings.TrimPrefix(authHeader[0], "Basic ")

	// Decode from Base64
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return
	}

	// Convert bytes to string
	decodedStr := string(decodedBytes)

	// Split into alias and public
	parts := strings.SplitN(decodedStr, ":", 2)
	if len(parts) != 2 {
		err = errors.New("Invalid Authorization Header")
		return
	}

	alias, password = parts[0], parts[1]
	return
}

func HashWithSalt(content, salt string) string {
	data := append([]byte(content), []byte(salt)...)
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

func createSession(alias string) (sessionId string, err error) {
	sessionId = uuid.New().String()
	err = insertSession(sessionId, alias)

	return
}

func setAuthCookie(w http.ResponseWriter, auth_token string) {
	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:     AUTH_COOKIE_NAME,
		Value:    auth_token,
		Expires:  expiration,
		HttpOnly: true,
		Secure:   true, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/", // Ensure the cookie is valid for all paths
	}
	http.SetCookie(w, &cookie)
}

func initiateSession(w http.ResponseWriter, alias string) {
	id, err := createSession(alias)
	if err != nil {
		panic(err) // TODO: Gracefully handle this error
	}
	log.Printf("Intiating session %s for alias %s\n", id, alias)
	setAuthCookie(w, id)
}
