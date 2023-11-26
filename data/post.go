package data

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	PublishedOn *time.Time `json:"published_on"`
	CreatedOn   time.Time  `json:"created_on"`
}

type Posts []Post

func (p *Post) Create(db *sql.DB) error {
	p.ID = uuid.NewString()
	stmt, err := db.Prepare(
		`
		INSERT INTO post (id, title, content) values($1, $2, $3) RETURNING created_on
		`,
	)
	if err != nil {
		return err
	}
	return stmt.QueryRow(p.ID, p.Title, p.Content).Scan(&p.CreatedOn)
}

func (p *Post) Publish(db *sql.DB) error {
	publishedOn := time.Now().UTC()
	p.PublishedOn = &publishedOn
	stmt, err := db.Prepare(
		`
		UPDATE post
		SET published_on = $1
		WHERE id = $2
		`,
	)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(p.PublishedOn, p.ID)
	return err
}

func (posts *Posts) ReadAllPosts(db *sql.DB) error {
	stmt, err := db.Prepare(
		`
		SELECT id, title, content, published_on, created_on
		FROM post
		`,
	)
	if err != nil {
		return err
	}

	rows, err := stmt.Query()
	if err != nil {
		return err
	}

	for rows.Next() {
		p := Post{}

		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.PublishedOn, &p.CreatedOn); err != nil {
			return err
		}

		*posts = append(*posts, p)
	}

	return nil
}
