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
	"github.com/tomihaapalainen/blog-api/schema"

	_ "github.com/mattn/go-sqlite3"
)

func TestPostNewPost(t *testing.T) {
	title := "Test Post"
	content := "Test post content."

	db, err := sql.Open("sqlite3", "file:../db.sqlite3?_fk=ON")
	if err != nil {
		t.Fatalf("err opening db: %+v", err)
	}
	e := echo.New()

	jsonStr := []byte(fmt.Sprintf(`{"title": "%s", "content": "%s"}`, title, content))
	req := httptest.NewRequest("POST", "/posts", bytes.NewBuffer(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	HandlePostPosts(db)(c)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected return code %d, was %d instead", http.StatusCreated, rec.Code)
	}
	post := model.Post{}
	if err := json.NewDecoder(rec.Body).Decode(&post); err != nil {
		t.Fatalf("unable to parse response: %+v", err)
	}
	_, err = uuid.Parse(post.ID)
	if err != nil {
		t.Fatalf("'%s' was not a valid UUID", post.ID)
	}
	if post.Title != title {
		t.Fatalf("post.Title '%s' != '%s'", post.Title, title)
	}
	if post.Content != content {
		t.Fatalf("post.Content '%s' != '%s'", post.Content, content)
	}
	if post.PublishedOn != nil {
		t.Fatalf("post.PublishedOn '%s' != nil", post.PublishedOn)
	}
}

func TestPostNewPostWithoutTitleShouldReturnBadRequest(t *testing.T) {
	content := "Test post content."

	db, err := sql.Open("sqlite3", "file:../db.sqlite3?_fk=ON")
	if err != nil {
		t.Fatalf("err opening db: %+v", err)
	}
	e := echo.New()

	jsonStr := []byte(fmt.Sprintf(`{"content": "%s"}`, content))
	req := httptest.NewRequest("POST", "/posts", bytes.NewBuffer(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	HandlePostPosts(db)(c)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected return code %d, was %d instead", http.StatusBadRequest, rec.Code)
	}
}

func TestPostNewPostWithoutContentShouldReturnBadRequest(t *testing.T) {
	title := "Test Post"

	db, err := sql.Open("sqlite3", "file:../db.sqlite3?_fk=ON")
	if err != nil {
		t.Fatalf("err opening db: %+v", err)
	}
	e := echo.New()

	jsonStr := []byte(fmt.Sprintf(`{"title": "%s"}`, title))
	req := httptest.NewRequest("POST", "/posts", bytes.NewBuffer(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	HandlePostPosts(db)(c)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected return code %d, was %d instead", http.StatusBadRequest, rec.Code)
	}
	res := schema.ErrorResponse{}
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("unable to parse response: %+v", err)
	}
}

func TestPostNewPostWithNullTitleShouldReturnBadRequest(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:../db.sqlite3?_fk=ON")
	if err != nil {
		t.Fatalf("err opening db: %+v", err)
	}
	e := echo.New()

	jsonStr := []byte(`{"title": null, "content": "Test post content."}`)
	req := httptest.NewRequest("POST", "/posts", bytes.NewBuffer(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	HandlePostPosts(db)(c)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected return code %d, was %d instead", http.StatusBadRequest, rec.Code)
	}
	res := schema.ErrorResponse{}
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("unable to parse response: %+v", err)
	}
}

func TestPostNewPostWithNullContentShouldReturnBadRequest(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:../db.sqlite3?_fk=ON")
	if err != nil {
		t.Fatalf("err opening db: %+v", err)
	}
	e := echo.New()

	jsonStr := []byte(`{"title": "Test Post", "content": null}`)
	req := httptest.NewRequest("POST", "/posts", bytes.NewBuffer(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	HandlePostPosts(db)(c)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected return code %d, was %d instead", http.StatusBadRequest, rec.Code)
	}
	res := schema.ErrorResponse{}
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("unable to parse response: %+v", err)
	}
}

func TestGetAllPosts(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:../db.sqlite3?_fk=ON")
	if err != nil {
		t.Fatalf("err opening db: %+v", err)
	}
	e := echo.New()

	req := httptest.NewRequest("GET", "/posts", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	HandleGetAllPosts(db)(c)
	if rec.Code != http.StatusOK {
		t.Fatalf("return code %d != %d", rec.Code, http.StatusOK)
	}
	res := model.Posts{}
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("unable to parse response: %+v", err)
	}
}
