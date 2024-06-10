package routes

import (
	"encoding/json"
	"forum/database"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// GetCategories retrieves all category names from the database
func GetCategories(w http.ResponseWriter, r *http.Request) {
	// Query the database to select all names from the categories table
	rows, err := database.DB.Query("SELECT name FROM categories")
	if err != nil {
		// If there is an error executing the query, return an internal server error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close() // Close the rows when the function returns

	var categories []string // Slice to store category names
	for rows.Next() {       // Iterate over the rows returned from the query
		var category string
		err := rows.Scan(&category) // Scan the current row into the category variable
		if err != nil {
			// If there is an error scanning the row, return an internal server error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		categories = append(categories, category) // Append the category name to the slice
	}

	// Set the response header to indicate JSON content
	w.Header().Set("Content-Type", "application/json")
	// Encode the categories slice as JSON and write it to the response
	json.NewEncoder(w).Encode(categories)
}
