package keyservice

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"zawie.io/e2e/backend/authservice"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Key service called", r.Method, r.URL.Path)

	if r.Method == "POST" {
		if registerKey(w, r) == nil {
			w.WriteHeader(http.StatusOK)
		}
		return
	}

	if r.Method == "GET" {
		subject := r.URL.Query().Get("subject")

		log.Println("Getting keys for", subject)
		keys := getKeys(subject)
		jsonData, err := json.Marshal(keys)

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

type registerKeyRequestData struct {
	DeviceUUID string `json:"deviceUUID"`
	KeyUUID    string `json:"keyUUID"`
	PublicKey  string `json:"publicKey"`
}

func registerKey(w http.ResponseWriter, r *http.Request) (err error) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Unmarshal the JSON data
	var data registerKeyRequestData
	err = json.Unmarshal(body, &data)

	alias := authservice.RetrieveAlias(r)
	log.Printf("Registration contents:\nUser\t%s\nDevice\t%s\nKey\t%s\n", alias, data.DeviceUUID, data.KeyUUID)

	if err != nil {
		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		return
	}

	insertKey(alias, data.DeviceUUID, data.KeyUUID, data.PublicKey)

	return
}
