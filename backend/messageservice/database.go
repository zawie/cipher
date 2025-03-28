package messageservice

import (
	"database/sql"
	"fmt"
	"log"
	"sort"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

var dbPath = "sample.db"

type Message struct {
	Sender    string   `json:"sender"`
	CreatedAt string   `json:"createdAt"`
	Recipient string   `json:"recipient"`
	Ciphers   []cipher `json:"ciphers"`
}

func queryMessages(sender, recipient string) (messages []Message) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	querySQL := `SELECT message_id, sender.alias, recipient.alias, message.created_at, key_id, message.cipher
FROM message
         JOIN user sender
              ON sender.id = message.sender_id
         JOIN user recipient
              ON recipient.id = message.recipient_id
WHERE (sender.alias = ? AND recipient.alias = ?)
   OR (recipient.alias = ? AND sender.alias = ?)`

	rows, err := db.Query(querySQL, sender, recipient, sender, recipient)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	msgMap := make(map[string]Message)
	for rows.Next() {
		var messageId string
		var sender string
		var recipient string
		var createdAt string
		var keyId string
		var ciphertext string

		err = rows.Scan(&messageId, &sender, &recipient, &createdAt, &keyId, &ciphertext)
		if err != nil {
			log.Fatal(err)
		}

		_, ok := msgMap[messageId]
		if !ok {
			msgMap[messageId] = Message{}
		}

		msg := msgMap[messageId]
		msg.Sender = sender
		msg.Recipient = recipient
		msg.CreatedAt = createdAt
		msg.Ciphers = append(msgMap[messageId].Ciphers, cipher{
			KeyUUID: keyId,
			Cipher:  ciphertext,
		})
		msgMap[messageId] = msg
	}

	for _, msg := range msgMap {
		messages = append(messages, msg)
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt > messages[j].CreatedAt
	})

	return
}

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
	err = db.QueryRow("SELECT id FROM user WHERE alias = ? LIMIT 1", alias).Scan(&id)

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
