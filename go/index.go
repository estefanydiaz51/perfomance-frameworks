package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Image struct {
	AlbumId      int    `json:"albumId"`
	Id           int    `json:"id"`
	Title        string `json:"title"`
	Url          string `json:"url"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

var db *sql.DB

func create(w http.ResponseWriter, r *http.Request) {

	fileContent, err := os.ReadFile("dummy.json")
	if err != nil {
		http.Error(w, "Error reading dummy.json", http.StatusInternalServerError)
		return
	}

	var images []Image
	if err := json.Unmarshal(fileContent, &images); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusInternalServerError)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Database transaction begin error", http.StatusInternalServerError)
		return
	}
	startTime := time.Now()
	for _, image := range images {
		_, err := tx.Exec("INSERT INTO images (albumId, id, title, url, thumbnailUrl) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING RETURNING *", image.AlbumId, image.Id, image.Title, image.Url, image.ThumbnailUrl)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error inserting data", http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Transaction commit error", http.StatusInternalServerError)
		return
	}

	endTime := time.Now()

	response := struct {
		Message      string `json:"message"`
		Time         string `json:"time"`
		DataQuantity int    `json:"dataQuantity"`
	}{
		Message:      "Datos cargados correctamente",
		Time:         fmt.Sprintf("%v ms", endTime.Sub(startTime).Milliseconds()),
		DataQuantity: len(images),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getImages(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	rows, err := db.Query("SELECT * FROM images ORDER BY id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var images []Image
	for rows.Next() {
		var img Image
		if err := rows.Scan(&img.AlbumId, &img.Id, &img.Title, &img.Url, &img.ThumbnailUrl); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		images = append(images, img)
	}
	endTime := time.Now()
	timeTaken := endTime.Sub(startTime)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":         images,
		"dataQuantity": len(images),
		"time":         fmt.Sprintf("%v ms", timeTaken.Milliseconds()),
	})
}

func updateImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var img Image
	if err := json.NewDecoder(r.Body).Decode(&img); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	_, err := db.Exec("UPDATE images SET albumId = $1, title = $2, url = $3, thumbnailUrl = $4 WHERE id = $5",
		img.AlbumId, img.Title, img.Url, img.ThumbnailUrl, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	endTime := time.Now()
	timeTaken := endTime.Sub(startTime)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Image updated successfully",
		"time":    fmt.Sprintf("%v ms", timeTaken.Milliseconds()),
	})
}

func deleteImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	startTime := time.Now()
	_, err := db.Exec("DELETE FROM images WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	endTime := time.Now()
	timeTaken := endTime.Sub(startTime)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Image deleted successfully",
		"time":    fmt.Sprintf("%v ms", timeTaken.Milliseconds()),
	})
}

func main() {
	var err error
	db, err = sql.Open("postgres", "postgres://backend:12345@localhost/crud_db?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r := mux.NewRouter()

	// Define your routes here
	r.HandleFunc("/create", create).Methods("POST")
	r.HandleFunc("/read", getImages).Methods("GET")
	r.HandleFunc("/update/{id}", updateImage).Methods("PUT")
	r.HandleFunc("/delete/{id}", deleteImage).Methods("DELETE")

	http.Handle("/", r)

	fmt.Println("Server running at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
