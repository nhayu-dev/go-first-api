package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite"
)

type message struct {
	Text string `json:"text"`
	ID   int    `json:"id"`
}

var db *sql.DB

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message{Text: "Hello, JSON API!"})
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	var msg message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if msg.Text == "" {
		http.Error(w, "text cannot be empty", http.StatusBadRequest)
		return
	}

	res, err := db.Exec("INSERT INTO messages (text) VALUES(?)", msg.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	msg.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)

}

func getHandler(w http.ResponseWriter, _ *http.Request) {
	rows, err := db.Query("SELECT id, text FROM messages")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var msgs []message
	for rows.Next() {
		var m message
		if err := rows.Scan(&m.ID, &m.Text); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		msgs = append(msgs, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msgs)
}

func getByIDHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Path[len("/messages/"):]
	var m message
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = db.QueryRow("SELECT id, text FROM messages WHERE id = ?", id).Scan(&m.ID, &m.Text)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "appliction/json")
	json.NewEncoder(w).Encode(m)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Path[len("/messages/"):]

	var newMsg message

	err := json.NewDecoder(r.Body).Decode(&newMsg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE messages SET text = ? WHERE id =?", newMsg.Text, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newMsg.ID = id

	json.NewEncoder(w).Encode(newMsg)

}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.URL.Path[len("/messages/"):])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	res, err := db.Exec("DELETE FROM messages WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("delited id %d", id),
	})

}

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./messages.db")
	if err != nil {
		panic(err)
	}

	createTable := `
		CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT NOT NULL
		);
		`

	_, err = db.Exec(createTable)
	if err != nil {
		panic(err)
	}
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Acesss-Control-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.DefaultServeMux.ServeHTTP(w, r)

	})

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			getHandler(w, r) //全件取得
		case http.MethodPost:
			postHandler(w, r) //新規作成
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/messages/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			getByIDHandler(w, r) //一件取得
		case http.MethodPut:
			updateHandler(w, r) //更新
		case http.MethodDelete:
			deleteHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		}
	})

	fmt.Println("Server running at http://localhost:3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
