package bbscommon

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func OpenDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTables(db *sql.DB) {
	createUserTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		address TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(createUserTable)
	if err != nil {
		log.Fatalf("Failed to create user table: %v", err)
	}

	createMessageBaseTable := `CREATE TABLE IF NOT EXISTS message_bases (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		access_read INTEGER,
		access_post INTEGER
	);`
	_, err = db.Exec(createMessageBaseTable)
	if err != nil {
		log.Fatalf("Failed to create message base table: %v", err)
	}

	createMessageTable := `CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		message_base_id INTEGER,
		content TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id),
		FOREIGN KEY(message_base_id) REFERENCES message_bases(id)
	);`
	_, err = db.Exec(createMessageTable)
	if err != nil {
		log.Fatalf("Failed to create message table: %v", err)
	}
}

func RegisterUser(db *sql.DB, name, password, address string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := "INSERT INTO users (name, password, address) VALUES (?, ?, ?)"
	_, err = db.Exec(query, name, string(hashedPassword), address)
	return err
}

func CheckPassword(db *sql.DB, userID int, password string) bool {
	var hashedPassword string
	query := "SELECT password FROM users WHERE id = ?"
	err := db.QueryRow(query, userID).Scan(&hashedPassword)
	if err != nil {
		log.Printf("Failed to retrieve password: %v", err)
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func GetUserID(db *sql.DB, name string) (int, error) {
	var userID int
	query := "SELECT id FROM users WHERE name = ?"
	err := db.QueryRow(query, name).Scan(&userID)
	return userID, err
}
