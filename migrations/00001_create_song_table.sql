-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS song (
    group_name VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL ,
    release_date DATE NOT NULL ,
    text TEXT NOT NULL ,
    link VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS song;
-- +goose StatementEnd
