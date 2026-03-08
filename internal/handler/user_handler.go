package handler

import (
	"encoding/json"
	"net/http"
	"project-start/internal/db"
	"project-start/internal/model"
	"strconv"
	"strings"
)
// ----- FUNCTION -----
// ----- CREATE -------
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// validate HTTP method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// deserialization (JSON -> Go struct)
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid json"})
		return
	}
	// parameterized SQL query (SQL Injections free)
	// numbered placeholders like $1 PostgreSQL
	// tells the driver exactly which argument to substitute
	query := `INSERT INTO users (email) VALUES ($1) RETURNING id, created_at`
	var id int
	var createdAt string
	// execute CREATE query
	err := db.DB.QueryRow(query, user.Email).Scan(&id, &createdAt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	// serialization
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         id,
		"email":      user.Email,
		"created_at": createdAt,
	})
}
// ----- FUNCTION -----
// ----- GET ----------
func GetUser(w http.ResponseWriter, r *http.Request) {
	// validate HTTP method
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// is id in query
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id is required"})
		return
	}
	// is id a number
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id must be a number"})
		return
	}
	// execute GET query
	var user model.User
	err = db.DB.QueryRow("SELECT id, email, created_at FROM users WHERE id=$1", id).Scan(
		&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	}
	// serialization
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ----- FUNCTION -----
// ----- DELETE -------
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// validate HTTP method
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// parse URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid path"})
		return
	}
	// is id a number
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id must be number"})
		return
	}
	// execute DELETE query
	result, err := db.DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	// check if user existed
	rows, _ := result.RowsAffected()
	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	}
	// serialization
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "user deleted",
	})

}

// ----- UPDATE -----
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// validate HTTP method
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// parse URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// is id a number
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// deserialization
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// execute UPDATE query
	result, err := db.DB.Exec(
    "UPDATE users SET email=$1 WHERE id=$2",
    user.Email,
    id,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// error handling
	rows, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	// serialization
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "user updated",
	})
}
