-- +goose Up
CREATE TABLE chat_rooms (
    id INT PRIMARY KEY NOT NULL,
    name text NOT NULL,
    timestamp TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE chat_rooms;