-- +goose Up
CREATE TABLE messages (
    id INT PRIMARY KEY NOT NULL,
    room_id INT NOT NULL,
    message TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    sent_from INT NOT NULL
);

-- +goose Down
DROP TABLE messages;