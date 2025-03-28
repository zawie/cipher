package keyservice

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

var dbPath = "sample.db"

func insertKey(alias, device, keyId, key string) {
	// Open (or create) an SQLite database file
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createKeyTable(db)

	// Insert data into the table
	insertSQL := fmt.Sprintf(`
	INSERT INTO key (user_id, device_id, key_id, key)
	SELECT  user.id, '%s', '%s', '%s'
	FROM    user
	WHERE   alias = '%s'
	`, device, keyId, key, alias)
	_, err = db.Exec(insertSQL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Data successfully inserted into %s\n", dbPath)
}

type Key struct {
	UUID    string
	Encoded string
}

func getKeys(alias string) (keys []Key) {
	// Open (or create) an SQLite database file
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createKeyTable(db)

	// Query the table
	querySQL := `SELECT key_id, key
FROM (SELECT user_id, key_id, key,
             RANK() OVER (
                 PARTITION BY key.user_id, key.device_id
                 ORDER BY created_at DESC
                 ) rank
      FROM key) t
JOIN user ON user.id = t.user_id
WHERE t.rank = 1 AND user.alias = ?;`

	rows, err := db.Query(querySQL, alias)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var key Key

		err = rows.Scan(&key.UUID, &key.Encoded)
		if err != nil {
			log.Fatal(err)
		}

		keys = append(keys, key)
	}

	fmt.Printf("Read %d keys for user %s from \"%s\"\n", len(keys), alias, dbPath)
	return
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
