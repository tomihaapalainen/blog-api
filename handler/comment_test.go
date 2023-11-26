package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/tomihaapalainen/blog-api/model"
)

var testPostID string

func setupSuite(tb testing.TB) {
	db, err := sql.Open("sqlite3", "file:../db.sqlite3?_fk=ON")
	if err != nil {
		tb.Fatalf("err opening db: %+v", err)
	}
	p := model.Post{Title: "Test Post", Content: "Test Content"}
	err = p.Create(db)
	if err != nil {
		tb.Fatalf("err creating post: %+v", err)
	}
	testPostID = p.ID
}

func TestPostComment(t *testing.T) {
	setupSuite(t)
	content := "Test comment content."

	db, err := sql.Open("sqlite3", "file:../db.sqlite3?_fk=ON")
	if err != nil {
		t.Fatalf("err opening db: %+v", err)
	}
	e := echo.New()

	jsonStr := []byte(fmt.Sprintf(`{"post_id": "%s", "content": "%s"}`, testPostID, content))
	req := httptest.NewRequest("POST", "/posts/comments", bytes.NewBuffer(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	HandlePostComment(db)(c)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected return code %d, was %d instead", http.StatusCreated, rec.Code)
	}
	comment := model.Comment{}
	if err := json.NewDecoder(rec.Body).Decode(&comment); err != nil {
		t.Fatalf("unable to parse response: %+v", err)
	}
	if _, err := uuid.Parse(comment.ID); err != nil {
		t.Fatalf("err parsing comment ID: %+v", err)
	}
	if comment.Content != content {
		t.Fatalf("content '%s' != '%s'", comment.Content, content)
	}
	if comment.PostID != testPostID {
		t.Fatalf("post ID '%s' != '%s'", comment.PostID, testPostID)
	}
}

func TestPostCommentWithoutPostIDShouldFail(t *testing.T) {
	setupSuite(t)

	content := "Test comment content."

	db, err := sql.Open("sqlite3", "file:../db.sqlite3?_fk=ON")
	if err != nil {
		t.Fatalf("err opening db: %+v", err)
	}
	e := echo.New()

	jsonStr := []byte(fmt.Sprintf(`{"content": "%s"}`, content))
	req := httptest.NewRequest("POST", "/posts/comments", bytes.NewBuffer(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	HandlePostComment(db)(c)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected return code %d, was %d instead", http.StatusCreated, rec.Code)
	}
}

func TestPostCommentWithoutContentShouldFail(t *testing.T) {
	setupSuite(t)

	db, err := sql.Open("sqlite3", "file:../db.sqlite3?_fk=ON")
	if err != nil {
		t.Fatalf("err opening db: %+v", err)
	}
	e := echo.New()

	jsonStr := []byte(fmt.Sprintf(`{"post_id": "%s"}`, testPostID))
	req := httptest.NewRequest("POST", "/posts/comments", bytes.NewBuffer(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	HandlePostComment(db)(c)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected return code %d, was %d instead", http.StatusCreated, rec.Code)
	}
}

func TestPostCommentWithWhiteSpaceContentShouldFail(t *testing.T) {
	setupSuite(t)

	db, err := sql.Open("sqlite3", "file:../db.sqlite3?_fk=ON")
	if err != nil {
		t.Fatalf("err opening db: %+v", err)
	}
	e := echo.New()

	jsonStr := []byte(fmt.Sprintf(`{"post_id": "%s", "content": "    "}`, testPostID))
	req := httptest.NewRequest("POST", "/posts/comments", bytes.NewBuffer(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	HandlePostComment(db)(c)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected return code %d, was %d instead", http.StatusCreated, rec.Code)
	}
}
