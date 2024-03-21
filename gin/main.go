package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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

func create(c *gin.Context) {
	fileContent, err := os.ReadFile("dummy.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading dummy.json"})
		return
	}

	var images []Image
	if err := json.Unmarshal(fileContent, &images); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing JSON"})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database transaction begin error"})
		return
	}
	startTime := time.Now()
	for _, image := range images {
		_, err := tx.Exec("INSERT INTO images (albumId, id, title, url, thumbnailUrl) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING RETURNING *", image.AlbumId, image.Id, image.Title, image.Url, image.ThumbnailUrl)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting data"})
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit error"})
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

	c.JSON(http.StatusOK, response)
}

func getImages(c *gin.Context) {
	startTime := time.Now()
	rows, err := db.Query("SELECT * FROM images ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var images []Image
	for rows.Next() {
		var img Image
		if err := rows.Scan(&img.AlbumId, &img.Id, &img.Title, &img.Url, &img.ThumbnailUrl); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		images = append(images, img)
	}
	endTime := time.Now()
	timeTaken := endTime.Sub(startTime)

	c.JSON(http.StatusOK, gin.H{
		"data":         images,
		"dataQuantity": len(images),
		"time":         fmt.Sprintf("%v ms", timeTaken.Milliseconds()),
	})
}

func updateImage(c *gin.Context) {
	id := c.Param("id")

	var img Image
	if err := c.BindJSON(&img); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startTime := time.Now()
	_, err := db.Exec("UPDATE images SET albumId = $1, title = $2, url = $3, thumbnailUrl = $4 WHERE id = $5",
		img.AlbumId, img.Title, img.Url, img.ThumbnailUrl, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	endTime := time.Now()
	timeTaken := endTime.Sub(startTime)

	c.JSON(http.StatusOK, gin.H{
		"message": "Image updated successfully",
		"time":    fmt.Sprintf("%v ms", timeTaken.Milliseconds()),
	})
}

func deleteImage(c *gin.Context) {
	id := c.Param("id")

	startTime := time.Now()
	_, err := db.Exec("DELETE FROM images WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	endTime := time.Now()
	timeTaken := endTime.Sub(startTime)

	c.JSON(http.StatusOK, gin.H{
		"message": "Image deleted successfully",
		"time":    fmt.Sprintf("%v ms", timeTaken.Milliseconds()),
	})
}

func main() {
	var err error
	db, err = sql.Open("postgres", "postgres://usuario:5131@localhost:5433/crud_db?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r := gin.Default()

	// Define your routes here
	r.POST("/create", create)
	r.GET("/read", getImages)
	r.PUT("/update/:id", updateImage)
	r.DELETE("/delete/:id", deleteImage)

	fmt.Println("Server running at http://52.23.103.132:3000")
	log.Fatal(r.Run("http://52.23.103.132:3000"))
}
