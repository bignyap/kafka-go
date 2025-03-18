-- +goose Up
CREATE TABLE room_members (
    room_id INT NOT NULL,
    member_id INT NOT NULL,
    PRIMARY KEY (room_id, member_id),
    FOREIGN KEY (room_id) REFERENCES chat_rooms(id),
    FOREIGN KEY (member_id) REFERENCES members(id)
);

-- +goose Down
DROP TABLE room_members;