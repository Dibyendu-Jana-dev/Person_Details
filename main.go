package main

import (
	"person_details/db"
	"person_details/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// here I Initialize the database connection
	db.InitDB()
	defer db.DB.Close()

	router := gin.Default()

	// Define the endpoint and handler
	router.GET("/person/:person_id/info", handlers.GetPersonInfo)
	router.POST("/person/create", handlers.CreatePerson)

	router.Run(":8080")
}
