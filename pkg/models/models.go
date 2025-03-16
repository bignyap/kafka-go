package models

import "time"

type Member struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	JoinedAt time.Time `json:"joined_at"`
}

type ChatMessage struct {
	ID        int       `json:"id"`
	RoomID    int       `json:"room_id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type ChatRoom struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Members   []int     `json:"members"`
	CreatedAt time.Time `json:"created_at"`
}
