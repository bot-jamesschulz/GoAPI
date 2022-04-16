package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model

	Title    string  `gorm:"uniqueIndex" json:"title"`
	Price    float32 `json:"price"`
	Director string  `json:"director"`
	Length   uint    `json:"length"`
	Rating   string  `json:"rating"`
}

var db *gorm.DB
var err error

func main() {
	// Env vars
	host := goDotEnvVariable("HOST")
	dbPort := goDotEnvVariable("DBPORT")
	user := goDotEnvVariable("USER")
	dbName := goDotEnvVariable("NAME")
	password := goDotEnvVariable("PASSWORD")

	// Connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbName, password, dbPort)

	// Connection
	db, err = gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successful connection")
	}

	// API routes
	router := gin.Default()
	router.GET("/movies", getMovies)
	router.POST("/movies", insertMovie)
	router.GET("/movies/:rating", getMoviesWithRating)
	router.DELETE("/movies/:id", deleteMovie)

	router.Run("localhost:8080")
}

// Returns all movies in the database
func getMovies(c *gin.Context) {
	var movies []Movie
	result := db.Find(&movies)
	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, result.Error)
	} else {
		c.IndentedJSON(http.StatusOK, &movies)
	}
}

// Adds a single movie to the database
func insertMovie(c *gin.Context) {
	var request Movie

	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	}

	result := db.Create(&request)
	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, result.Error)
	} else {
		c.IndentedJSON(http.StatusOK, &request)
	}

}

// Gets all movies with specified rating
func getMoviesWithRating(c *gin.Context) {
	rating := c.Param("rating")
	var movies []Movie

	result := db.Where("rating = ?", rating).Find(&movies)
	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, result.Error)
	} else {
		c.IndentedJSON(http.StatusOK, &movies)
	}

}

// Deletes movie with specified id
func deleteMovie(c *gin.Context) {
	id := c.Param("id")

	result := db.Delete(&Movie{}, id)
	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, result.Error)
	} else {
		c.Status(200)
	}

}

// Get env variables from .env
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
