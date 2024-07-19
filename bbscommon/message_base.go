package bbscommon

import (
	"database/sql"
	"fmt"
)

type MessageBase struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	AccessRead int    `json:"access_read"`
	AccessPost int    `json:"access_post"`
}

func CreateMessageBase(db *sql.DB, name string, accessRead, accessPost int) error {
	query := `INSERT INTO message_bases (name, access_read, access_post) VALUES (?, ?, ?)`
	_, err := db.Exec(query, name, accessRead, accessPost)
	return err
}

func UpdateMessageBase(db *sql.DB, id int, name string, accessRead, accessPost int) error {
	query := `UPDATE message_bases SET name = ?, access_read = ?, access_post = ? WHERE id = ?`
	_, err := db.Exec(query, name, accessRead, accessPost, id)
	return err
}

func DeleteMessageBase(db *sql.DB, id int) error {
	query := `DELETE FROM message_bases WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

func ListMessageBases(db *sql.DB) ([]MessageBase, error) {
	query := `SELECT id, name, access_read, access_post FROM message_bases`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messageBases []MessageBase
	for rows.Next() {
		var mb MessageBase
		if err := rows.Scan(&mb.ID, &mb.Name, &mb.AccessRead, &mb.AccessPost); err != nil {
			return nil, err
		}
		messageBases = append(messageBases, mb)
	}

	return messageBases, nil
}

func GetMessageBaseID(db *sql.DB, name string) (int, error) {
	var id int
	query := "SELECT id FROM message_bases WHERE name = ?"
	err := db.QueryRow(query, name).Scan(&id)
	return id, err
}

func ViewMessages(db *sql.DB, messageBaseID int) ([]string, error) {
	query := `SELECT users.name, messages.content, messages.created_at 
			  FROM messages 
			  JOIN users ON messages.user_id = users.id 
			  WHERE messages.message_base_id = ?
			  ORDER BY messages.created_at DESC`
	rows, err := db.Query(query, messageBaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []string
	for rows.Next() {
		var name, content string
		var createdAt string
		err := rows.Scan(&name, &content, &createdAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, fmt.Sprintf("%s (%s): %s", name, createdAt, content))
	}

	return messages, nil
}
