package messageservice

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

var dbPath = "sample.db"

func insertMessage(message_id, sender, recipient string, ciphers []cipher) (err error) {
	// Open (or create) an SQLite database file
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	creatMessageTable(db)
	
	sender_id, err := getUserIdFromAlias(db, sender)
	if err != nil {
		log.Printf("Could not get id for sender %s\n", sender)
		return
	}
	recipient_id, err := getUserIdFromAlias(db, recipient)
	if err != nil {
		log.Printf("Could not get id for recipoent %s\n", recipient)
		return
	}
	
	// Insert data into the table
	insertSQL := `INSERT INTO message (message_id, recipient_id, sender_id, key_id, cipher)	VALUES`

	for i, cipher := range ciphers {
		if i > 0 {
			insertSQL += ","
		}
		insertSQL += fmt.Sprintf("\n('%s', '%s', '%s', '%s', '%s')", message_id, recipient_id, sender_id, cipher.KeyUUID, cipher.Cipher)
	}
	_, err = db.Exec(insertSQL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Message successfully (%d ciphers) inserted into %s\n", len(ciphers), dbPath)
	
	return
}

func getUserIdFromAlias(db *sql.DB, alias string) (id string, err error) {
	err = db.QueryRow("SELECT user_id FROM user WHERE alias = ? LIMIT 1", alias).Scan(&id)
	
	return
}

func creatMessageTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS message (
		message_id 		TEXT NOT NULL,
		recipient_id	TEXT NOT NULL,
		sender_id		TEXT NOT NULL,
		key_id			TEXT NOT NULL,
		cipher			TEXT NOT NULL,
		created_at		DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}
