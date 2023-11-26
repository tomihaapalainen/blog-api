-- +goose Up
-- +goose StatementBegin
CREATE TABLE post (
    id TEXT PRIMARY KEY,
    title TEXT,
    content TEXT,
    published_on TIMESTAMP DEFAULT NULL,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE comment (
    id TEXT PRIMARY KEY,
    post_id TEXT,
    content TEXT,
    FOREIGN KEY(post_id) REFERENCES post(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE comment;
DROP TABLE post;
-- +goose StatementEnd
