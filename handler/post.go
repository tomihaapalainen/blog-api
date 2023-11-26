package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/tomihaapalainen/blog-api/data"
)

func HandlePostPosts(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		p := data.Post{}
		if err := json.NewDecoder(c.Request().Body).Decode(&p); err != nil {
			log.Println(fmt.Sprintf("Invalid post data: %+v", err))
			return c.JSON(
				http.StatusBadRequest,
				data.ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Message:    fmt.Sprintf("Invalid post data: %+v", err),
				},
			)
		}

		p.Title = strings.TrimSpace(p.Title)
		p.Content = strings.TrimSpace(p.Content)

		if p.Title == "" || p.Content == "" {
			return c.JSON(
				http.StatusBadRequest,
				data.ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Message:    "Title and content must not be empty",
				},
			)
		}

		if err := p.Create(db); err != nil {
			log.Println(fmt.Sprintf("Unable to create new post: %+v", err))
			return c.JSON(
				http.StatusBadRequest,
				data.ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Message:    fmt.Sprintf("Unable to create new post: %+v", err),
				},
			)
		}
		return c.JSON(http.StatusCreated, p)
	})
}

func HandleGetAllPosts(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		posts := data.Posts{}
		if err := posts.ReadAllPosts(db); err != nil {
			log.Println(fmt.Sprintf("Error reading all posts: %+v", err))
			return c.JSON(
				http.StatusInternalServerError,
				data.ErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    fmt.Sprintf(fmt.Sprintf("Error reading all posts: %+v", err)),
				},
			)
		}
		return c.JSON(http.StatusOK, posts)
	})
}
