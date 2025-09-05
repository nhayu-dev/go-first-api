package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type message struct {
	Text string `json:"text"`
	ID   int    `json:"id"`
}

var nextID = 1
var messages []message

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message{Text: "Hello, JSON API!"})
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	var msg message

	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg.ID = nextID
	nextID++

	messages = append(messages, msg)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)

}

func getHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func getByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "apploction/json")

	idStr := r.URL.Path[len("/get/"):]
	var id int
	fmt.Sscanf(idStr, "%d", &id)

	for _, m := range messages {
		if m.ID == id {
			json.NewEncoder(w).Encode(m)
			return
		}
	}

	http.Error(w, "notfound", http.StatusNotFound)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Path[len("/get/"):]
	var id int
	fmt.Sscanf(idStr, "%d", &id)

	var newMsg message

	err := json.NewDecoder(r.Body).Decode(&newMsg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, m := range messages {
		if m.ID == id {
			messages[i].Text = newMsg.Text
			json.NewEncoder(w).Encode(messages[i])
			return
		}
	}
	http.Error(w, "not found", http.StatusNotFound)

}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/get/", getByIDHandler)
	http.HandleFunc("/get/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getByIDHandler(w, r)
		case http.MethodPut:
			updateHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		}
	})

	fmt.Println("Server running at http://localhost:3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Eroor:", err)
	}
}
