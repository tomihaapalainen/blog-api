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
	"github.com/tomihaapalainen/blog-api/data"
)

func setupTest(tb testing.TB) (*sql.DB, string) {
	db, err := sql.Open("sqlite3", "file:../db.sqlite3?_fk=ON")
	if err != nil {
		tb.Fatalf("err opening db: %+v", err)
	}
	p := data.Post{Title: "Test Post", Content: "Test Content"}
	err = p.Create(db)
	if err != nil {
		tb.Fatalf("err creating post: %+v", err)
	}
	c := data.Comment{PostID: p.ID, Content: "Test Comment"}
	if err := c.Create(db); err != nil {
		tb.Fatalf("err creating comment: %+v", err)
	}

	return db, p.ID
}

func TestPostComment(t *testing.T) {
	db, testPostID := setupTest(t)
	content := "Test comment content."

	e := echo.New()

	jsonStr := []byte(fmt.Sprintf(`{"post_id": "%s", "content": "%s"}`, testPostID, content))
	req := httptest.NewRequest("POST", "/posts/comments", bytes.NewBuffer(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	HandlePostComment(db)(c)
	if rec.Code != http.StatusCreated {
		t.Fatalf("return code %d != %d", rec.Code, http.StatusCreated)
	}
	comment := data.Comment{}
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
	db, _ := setupTest(t)
	content := "Test comment content."

	e := echo.New()

	jsonStr := []byte(fmt.Sprintf(`{"content": "%s"}`, content))
	req := httptest.NewRequest("POST", "/posts/comments", bytes.NewBuffer(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	HandlePostComment(db)(c)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("return code %d != %d", rec.Code, http.StatusCreated)
	}
}

func TestPostCommentWithoutContentShouldFail(t *testing.T) {
	db, testPostID := setupTest(t)

	e := echo.New()

	jsonStr := []byte(fmt.Sprintf(`{"post_id": "%s"}`, testPostID))
	req := httptest.NewRequest("POST", "/posts/comments", bytes.NewBuffer(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	HandlePostComment(db)(c)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("return code %d != %d", rec.Code, http.StatusCreated)
	}
}

func TestPostCommentWithWhiteSpaceContentShouldFail(t *testing.T) {
	db, testPostID := setupTest(t)

	e := echo.New()

	jsonStr := []byte(fmt.Sprintf(`{"post_id": "%s", "content": "    "}`, testPostID))
	req := httptest.NewRequest("POST", "/posts/comments", bytes.NewBuffer(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	HandlePostComment(db)(c)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("return code %d != %d", rec.Code, http.StatusCreated)
	}
}

func TestGetPostComments(t *testing.T) {
	db, testPostID := setupTest(t)

	e := echo.New()

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/posts/:id/comments")
	c.SetParamNames("id")
	c.SetParamValues(testPostID)

	HandleGetPostComments(db)(c)
	if rec.Code != http.StatusOK {
		t.Fatalf("return code %d != %d", rec.Code, http.StatusOK)
	}

	comments := data.Comments{}
	if err := json.NewDecoder(rec.Body).Decode(&comments); err != nil {
		t.Fatalf("unable to parse json: %+v", err)
	}
}
