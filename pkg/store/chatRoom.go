package store

import (
	"database/sql"

	"github.com/bignyap/kafka-go/pkg/models"
)

type ChatRoomStore struct {
	db *sql.DB
}

func (chatRoomStore *ChatRoomStore) AddMemberToRoom(roomID int, memberID int) error {
	_, err := chatRoomStore.db.Exec(
		"INSERT INTO room_members (room_id, member_id) VALUES (?, ?)",
		roomID, memberID,
	)
	return err
}

func (chatRoomStore *ChatRoomStore) RemoveMemberFromRoom(roomID int, memberID int) error {
	_, err := chatRoomStore.db.Exec(
		"DELETE FROM room_members WHERE room_id = ? AND member_id = ?",
		roomID, memberID,
	)
	return err
}

func (chatRoomStore *ChatRoomStore) GetMembersFromRoom(DB *sql.DB, roomID int) ([]models.Member, error) {
	rows, err := chatRoomStore.db.Query(
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

func (chatRoomStore *ChatRoomStore) GetChatRoomsForMember(DB *sql.DB, memberID int) ([]models.ChatRoom, error) {
	rows, err := chatRoomStore.db.Query(
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
