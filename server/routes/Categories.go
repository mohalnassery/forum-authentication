package routes

import (
	"encoding/json"
	"forum/database"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT name FROM categories")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
