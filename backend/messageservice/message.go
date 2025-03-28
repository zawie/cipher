package messageservice

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"zawie.io/e2e/backend/authservice"
)

type postMessageRequest struct {
	Recipient string   `json:"recipient"`
	Ciphers   []cipher `json:"ciphers"`
}

type cipher struct {
	Cipher  string `json:"cipher"`
	KeyUUID string `json:"keyUUID"`
}

type getMessageResponse struct {
	Messages []Message `json:"messages"`
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

		sender := authservice.RetrieveAlias(r)
		recipient := data.Recipient

		log.Printf("Inserting messages from \"%s\" to \"%s\" into the database\n", sender, recipient)
		message_id := uuid.New().String()
		insertMessage(message_id, sender, recipient, data.Ciphers)

		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == "GET" {
		sender := authservice.RetrieveAlias(r)
		recipient := r.URL.Query().Get("subject")
		if recipient == "" {
			http.Error(w, "Subject query parameter is required", http.StatusBadRequest)
			return
		}

		log.Printf("Getting messages between \"%s\" and \"%s\"", sender, recipient)
		messages := queryMessages(sender, recipient)

		log.Printf("Sending messages: %v\n", messages)
		jsonData, err := json.Marshal(getMessageResponse{
			Messages: messages,
		})

		if err != nil {
			http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
