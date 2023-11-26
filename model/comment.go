package model

import (
	"database/sql"

	"github.com/google/uuid"
)

type Comment struct {
	ID      string `json:"id"`
	PostID  string `json:"post_id"`
	Content string `json:"content"`
}

func (c *Comment) Create(db *sql.DB) error {
	c.ID = uuid.NewString()
	stmt, err := db.Prepare(
		`
		INSERT INTO comment (id, post_id, content) values($1, $2, $3)
		`,
	)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(c.ID, c.PostID, c.Content)
	return err
}
