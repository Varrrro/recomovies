package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/varrrro/recomovies/internal/database"
	"github.com/varrrro/recomovies/internal/handler"
)

func main() {
	// Create database file if it doesn't exist
	create := false
	if _, err := os.Stat("db.sqlite3"); os.IsNotExist(err) {
		os.Create("db.sqlite3")
		create = true
	}

	// Open database connection and create data manager
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatalf("Can't open database connection, error: %s\n", err.Error())
	}
	defer db.Close()
	mg := database.NewManager(db)

	// Create and populate database if needed
	if create {
		mg.UpdateSchemas()
		if err := database.Populate(mg); err != nil {
			log.Fatalf("Can't populate database from files, error: %s\n", err.Error())
		}
	}

	// Set up web server
	h := handler.New(mg)
	r := gin.Default()
	r.LoadHTMLGlob("web/templates/*")

	// Set up routes
	r.GET("/", h.GetIndex)
	r.POST("/", h.PostIndex)
	r.GET("/rating/:n", h.GetRating)
	r.POST("/rating/:n", h.PostRating)
	r.GET("/results", h.GetResults)

	// Run web server
	r.Run(":8000")
}
