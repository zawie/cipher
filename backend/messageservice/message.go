package messageservice

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type postMessageRequest struct {
	Recipient string   `json:"recipient"`
	Ciphers   []cipher `json:"ciphers"`
}

type cipher struct {
	Cipher  string `json:"cipher"`
	KeyUUID string `json:"keyUUID"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Message service called", r.Method, r.URL.Path)

	if r.Method == "POST" {
		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Unmarshal the JSON data
		var data postMessageRequest
		err = json.Unmarshal(body, &data)
		
		if err != nil {
			log.Println(err)
			http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
			return
		}

		//TODO: Lol don't hardcode this
		sender := "anya"
		if data.Recipient == "anya" {
			sender = "zawie"
		}

		recipient := data.Recipient

		log.Println("Inserting messages to database")
		message_id := uuid.New().String()
		insertMessage(message_id, sender, recipient, data.Ciphers)

		w.WriteHeader(http.StatusOK)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
