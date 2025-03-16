package models

import "database/sql"

type ChatRoomManager interface {
	AddMember(roomID int, memberID int) error
	RemoveMember(roomID int, memberID int) error
	GetMembers(roomID int) ([]Member, error)
	GetChatRooms(memberID int) ([]ChatRoom, error)
}

type ChatMessageManager interface {
	SendMessage(roomID int, message string) error
	GetMessages(roomID int) ([]ChatMessage, error)
}

type ChatRoomManagerImpl struct {
	DB *sql.DB
}

func (crm *ChatRoomManagerImpl) AddMember(roomID int, memberID int) error {
	_, err := crm.DB.Exec(
		"INSERT INTO room_members (room_id, member_id) VALUES (?, ?)",
		roomID, memberID,
	)
	return err
}

func (crm *ChatRoomManagerImpl) RemoveMember(roomID int, memberID int) error {
	_, err := crm.DB.Exec(
		"DELETE FROM room_members WHERE room_id = ? AND member_id = ?",
		roomID, memberID,
	)
	return err
}

func (crm *ChatRoomManagerImpl) GetMembers(roomID int) ([]Member, error) {
	rows, err := crm.DB.Query(
		"SELECT member_id FROM room_members WHERE room_id = ?",
		roomID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var member Member
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

func (crm *ChatRoomManagerImpl) GetChatRooms(memberID int) ([]ChatRoom, error) {
	rows, err := crm.DB.Query(
		"SELECT room_id FROM room_members WHERE member_id = ?",
		memberID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []ChatRoom
	for rows.Next() {
		var room ChatRoom
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

type ChatMessageManagerImpl struct {
	DB *sql.DB
}

func (cmm *ChatMessageManagerImpl) SendMessage(roomID int, message string) error {
	_, err := cmm.DB.Exec(
		"INSERT INTO messages (room_id, message) VALUES (?, ?)",
		roomID, message,
	)
	return err
}

func (cmm *ChatMessageManagerImpl) GetMessages(roomID int) ([]ChatMessage, error) {
	rows, err := cmm.DB.Query(
		"SELECT id, message, timestamp FROM messages WHERE room_id = ?",
		roomID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var message ChatMessage
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
