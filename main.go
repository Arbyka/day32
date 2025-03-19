package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	users  = []User{{ID: 1, Name: "John Doe"}, {ID: 2, Name: "Jane Doe"}}
	mu     sync.Mutex
	nextID = 3
)

// getUsers handles GET /api/users
func getUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /api/users called") // Debug log
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Set timeout untuk memastikan respon cepat
	time.AfterFunc(2*time.Second, func() {
		json.NewEncoder(w).Encode(users)
	})
}

// addUser handles POST /api/users
func addUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /api/users called") // Debug log
	var newUser User

	// Decode dengan batas waktu agar tidak menggantung
	decodeDone := make(chan error, 1)
	go func() {
		decodeDone <- json.NewDecoder(r.Body).Decode(&newUser)
	}()

	select {
	case err := <-decodeDone:
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case <-time.After(2 * time.Second):
		http.Error(w, "Request Timeout", http.StatusRequestTimeout)
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

// Handler is the main entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler function triggered") // Debug log

	// Menggunakan multiplexer agar hanya rute /api/users diterima
	mux := http.NewServeMux()
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getUsers(w, r)
		case http.MethodPost:
			addUser(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Menjalankan multiplexer
	mux.ServeHTTP(w, r)
}
