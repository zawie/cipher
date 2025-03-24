package authservice

import (
	"log"
	"net/http"
)

var ALIAS_KEY string = "ALIAS"

func AuthMiddleware(redirect bool, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias, authenticated, err := checkAuthentication(r)

		log.Printf("Alias=%s; authenticated=%t\n", alias, authenticated)

		if err != nil {
			log.Printf("Something went wrong in the Authentication Middleware: %s\n", err.Error())
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		if !authenticated {
			if redirect {
				log.Println("Unauthenticated redirecting to /login!")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			} else {
				log.Println("Unauthenticated returning 401")
				http.Error(w, "Unauthenticated", http.StatusUnauthorized)
			}
			return
		}

		next(w, r)
	}
}
