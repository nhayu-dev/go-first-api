package main

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "modernc.org/sqlite"
)

func initTestDB(t *testing.T) *sql.DB {

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("faitaled to open test DB:%v", err)

	}

	_, err = db.Exec(`
CREATE TABLE messages(
id INTEGER PRIMARY KEY AUTOINCREMENT,
text TEXT NOT NULL);`)

	if err != nil {
		t.Fatalf("faitaled to open test DB:%v", err)
	}

	return db
}

func TestPostMessage_Success(t *testing.T) {
	db := initTestDB(t)
	defer db.Close()

	body := bytes.NewBufferString(`{"text":"hello test"}`)
	req := httptest.NewRequest(http.MethodPost, "/messages", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	postHandler(db, w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)

	}
}

func TestPostMessage_Fail_EmptyText(t *testing.T) {
	db := initTestDB(t)
	defer db.Close()

	body := bytes.NewBufferString(`{"text":""}`)
	req := httptest.NewRequest(http.MethodPost, "/messages", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	postHandler(db, w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 BadReqest, got %d", w.Code)
	}
}
