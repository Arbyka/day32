package handler

import (
	"encoding/json"
	"net/http"
	"sync"
)

// Struktur User
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	users  = []User{{ID: 1, Name: "John Doe"}, {ID: 2, Name: "Jane Doe"}}
	mu     sync.Mutex
	nextID = 3
)

// GET: Mengembalikan daftar pengguna
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// POST: Menambahkan pengguna baru
func addUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	newUser.ID = nextID
	nextID++
	users = append(users, newUser)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// Handler utama untuk menangani request
func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w, r)
	case http.MethodPost:
		addUser(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
