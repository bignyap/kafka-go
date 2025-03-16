package db

import (
	"database/sql"

	"github.com/bignyap/kafka-go/pkg/models"
)

func AddMemberToRoom(DB *sql.DB, roomID int, memberID int) error {
	_, err := DB.Exec(
		"INSERT INTO room_members (room_id, member_id) VALUES (?, ?)",
		roomID, memberID,
	)
	return err
}

func RemoveMemberFromRoom(DB *sql.DB, roomID int, memberID int) error {
	_, err := DB.Exec(
		"DELETE FROM room_members WHERE room_id = ? AND member_id = ?",
		roomID, memberID,
	)
	return err
}

func GetMembersFromRoom(DB *sql.DB, roomID int) ([]models.Member, error) {
	rows, err := DB.Query(
		"SELECT member_id FROM room_members WHERE room_id = ?",
		roomID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.Member
	for rows.Next() {
		var member models.Member
		if err := rows.Scan(&member.ID); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

func GetChatRoomsForMember(DB *sql.DB, memberID int) ([]models.ChatRoom, error) {
	rows, err := DB.Query(
		"SELECT room_id FROM room_members WHERE member_id = ?",
		memberID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.ChatRoom
	for rows.Next() {
		var room models.ChatRoom
		if err := rows.Scan(&room.ID); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

func SendMessageToRoom(DB *sql.DB, roomID int, message string) error {
	_, err := DB.Exec(
		"INSERT INTO messages (room_id, message) VALUES (?, ?)",
		roomID, message,
	)
	return err
}

func GetMessagesFromRoom(DB *sql.DB, roomID int) ([]models.ChatMessage, error) {
	rows, err := DB.Query(
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
