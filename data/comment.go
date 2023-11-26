package data

import (
	"database/sql"

	"github.com/google/uuid"
)

type Comment struct {
	ID      string `json:"id"`
	PostID  string `json:"post_id"`
	Content string `json:"content"`
}

type Comments []Comment

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

func (cs *Comments) ReadAllPostComments(db *sql.DB, postID string) error {
	stmt, err := db.Prepare(
		`
		SELECT id, post_id, content
		FROM comment
		WHERE post_id = $1
		`,
	)
	if err != nil {
		return err
	}

	rows, err := stmt.Query(postID)
	if err != nil {
		return err
	}

	for rows.Next() {
		c := Comment{}

		if err := rows.Scan(&c.ID, &c.PostID, &c.Content); err != nil {
			return err
		}

		*cs = append(*cs, c)
	}

	return nil
}
