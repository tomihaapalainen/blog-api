package main

import (
	"database/sql"
	"log"

	"github.com/labstack/echo"
	"github.com/tomihaapalainen/blog-api/handler"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "file:.//db.sqlite3?_fk=ON")
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.POST("/posts", handler.HandlePostPosts(db))
	e.GET("/posts", handler.HandleGetAllPosts(db))
	e.POST("/posts/comments", handler.HandlePostComment(db))
	e.GET("/posts/:id/comments", handler.HandleGetPostComments(db))

	e.Start(":8080")
}
