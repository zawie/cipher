package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

var dbPath = "sample.db"

func InsertKey(user, device, keyId, key string) {
	// Open (or create) an SQLite database file
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	createKeyTable(db)

	// Insert data into the table
	insertSQL := `INSERT INTO key (user_id, device_id, key_id, key) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(insertSQL, user, device, keyId, key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Data successfully inserted into %s\n", dbPath)
}

func createKeyTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS key (
		user_id 	TEXT NOT NULL,
		device_id	TEXT NOT NULL,
		key_id		TEXT NOT NULL,
		key			TEXT NOT NULL,
		created_at	DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}