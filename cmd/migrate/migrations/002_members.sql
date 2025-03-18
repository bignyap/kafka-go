-- +goose Up
CREATE TABLE members (
    id INT PRIMARY KEY NOT NULL,
    name text NOT NULL,
    email text NOT NULL,
    gender text NOT NULL,
    dob date NOT NULL,
    country text NOT NULL
);

-- +goose Down
DROP TABLE members;