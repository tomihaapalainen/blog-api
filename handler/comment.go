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

func HandlePostComment(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		comment := data.Comment{}
		if err := json.NewDecoder(c.Request().Body).Decode(&comment); err != nil {
			log.Println(fmt.Sprintf("err parsing json: %+v", err))
			return c.JSON(
				http.StatusBadRequest,
				data.ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Message:    fmt.Sprintf("err parsing json: %+v", err),
				},
			)
		}

		comment.PostID = strings.TrimSpace(comment.PostID)
		comment.Content = strings.TrimSpace(comment.Content)

		if comment.PostID == "" || comment.Content == "" {
			log.Println("Post id and content cannot be empty")
			return c.JSON(
				http.StatusBadRequest,
				data.ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Message:    "Post id and content cannot be empty",
				},
			)
		}

		if err := comment.Create(db); err != nil {
			log.Println(fmt.Sprintf("err creating comment: %+v", err))
			return c.JSON(
				http.StatusInternalServerError,
				data.ErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    fmt.Sprintf("err creating comment: %+v", err),
				},
			)
		}

		return c.JSON(http.StatusCreated, comment)
	})
}

func HandleGetPostComments(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		postID := c.Param("postID")
		comments := data.Comments{}
		if err := comments.ReadAllPostComments(db, postID); err != nil {
			log.Printf("err reading comments: %+v", err)
			return c.JSON(
				http.StatusInternalServerError,
				data.ErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    fmt.Sprintf("Error reading comments: %+v", err),
				},
			)
		}

		return c.JSON(http.StatusOK, comments)
	})
}
