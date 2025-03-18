package store

import (
	"database/sql"

	"github.com/bignyap/kafka-go/pkg/models"
)

type MessageStore struct {
	db *sql.DB
}

func (msgStore *MessageStore) SendMessageToRoom(roomID int, message string) error {
	_, err := msgStore.db.Exec(
		"INSERT INTO messages (room_id, message) VALUES (?, ?)",
		roomID, message,
	)
	return err
}

func (msgStore *MessageStore) GetMessagesFromRoom(roomID int) ([]models.ChatMessage, error) {
	rows, err := msgStore.db.Query(
		"SELECT id, message, timestamp FROM messages WHERE room_id = ?",
		roomID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.ChatMessage
	for rows.Next() {
		var message models.ChatMessage
		if err := rows.Scan(
			&message.ID, &message.Message, &message.Timestamp,
		); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
