package authservice

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

var dbPath = "sample.db"

func CheckIfAliasExists(alias string) (exists bool, err error) {
	// Open (or create) an SQLite database file
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createUserTable(db)

	query := "SELECT EXISTS(SELECT 1 FROM user WHERE alias = ?)"
	err = db.QueryRow(query, alias).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func insertUser(id, alias, salt, hashedAndSaltedPassword string) (err error) {
	// Open (or create) an SQLite database file
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	createUserTable(db)

	insertSQL := `INSERT INTO user (id, alias, salt, password) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(insertSQL, id, alias, salt, hashedAndSaltedPassword)
	if err != nil {
		log.Fatal(err)
		return
	}

	return
}

func pullSession(id string, maxLife time.Duration) (alias string, err error) {
	// Open (or create) an SQLite database file
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createSessionTable(db)

	query := fmt.Sprintf("SELECT alias FROM session WHERE session.id = ? AND created_at > datetime('now', '%f hours')", maxLife.Hours())
	err = db.QueryRow(query, id).Scan(&alias)

	return
}

func insertSession(id, alias string) (err error) {
	// Open (or create) an SQLite database file
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	createSessionTable(db)

	insertSQL := `INSERT INTO session (id, alias) VALUES (?, ?)`
	_, err = db.Exec(insertSQL, id, alias)
	if err != nil {
		log.Fatal(err)
		return
	}

	return
}

func createUserTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS user (
		id 			TEXT NOT NULL,
		alias		TEXT NOT NULL,
		salt		TEXT NOT NULL,
		password	TEXT NOT NULL,
		created_at	DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func createSessionTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS session (
		id 			TEXT NOT NULL,
		alias		TEXT NOT NULL,
		created_at	DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}
