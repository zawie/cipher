package keyservice

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	
	"zawie.io/e2e/backend/database"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Key service called", r.Method)
	
	if r.Method == "POST" {		
		if RegisterKey(w, r) == nil {
			w.WriteHeader(http.StatusOK)
		}
		return
	}
	
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

type RegisterKeyRequestData struct {
	DeviceUUID  string `json:"deviceUUID"`
	KeyUUID string    `json:"keyUUID"`
	PublicKey string `json:"publicKey"`
}

func RegisterKey(w http.ResponseWriter, r *http.Request) (err error) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Unmarshal the JSON data
	var data RegisterKeyRequestData
	err = json.Unmarshal(body, &data)
	
	// TODO: Pull user UUID from auth; for now just use device UUID
	UserUUID := data.DeviceUUID
	log.Printf("Registration contents:\nUser\t%s\nDevice\t%s\nKey\t%s\nKey\t%v", UserUUID, data.DeviceUUID, data.KeyUUID, data.PublicKey)
	
	if err != nil {
		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		return
	}
	
	database.InsertKey(UserUUID, data.DeviceUUID, data.KeyUUID, data.PublicKey)
	
	return
}


 