package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"person_details/constants"
	"person_details/db"
	"person_details/models"

	"github.com/gin-gonic/gin"
)

//create the person
func CreatePerson(c *gin.Context) {
	var person models.PersonInfo

	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	var personID int64
	query := "INSERT INTO person(name, age) VALUES (?, ?)"
	result, err := tx.Exec(query, person.Name, 10) // here i keep age static tusk1 and tusk2 didn't saying about age. but table schema's mentioned about the age
	if err != nil {
		log.Fatal(err)
	}
	personID, err = result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	query = "INSERT INTO phone(person_id, number) VALUES (?, ?)"
	_, err = tx.Exec(query, personID, person.PhoneNumber)
	if err != nil {
		log.Fatal(err)
	}

	query = "INSERT INTO address(city, state, street1, street2, zip_code) VALUES (?, ?, ?, ?, ?)"
	result, err = tx.Exec(query, person.City, person.State, person.Street1, person.Street2, person.ZipCode)
	if err != nil {
		log.Fatal(err)
	}
	addressID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	query = "INSERT INTO address_join(person_id, address_id) VALUES (?, ?)"
	_, err = tx.Exec(query, personID, addressID)
	if err != nil {
		log.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"status": constants.PERSON_CREATED_SUCCESSFULLY})
}
//fetching the person
func GetPersonInfo(c *gin.Context) {
	personID := c.Param("person_id")

	var personInfo models.PersonInfo

	query := `
		SELECT p.name, ph.number, a.city, a.state, a.street1, a.street2, a.zip_code
		FROM person p
		JOIN phone ph ON p.id = ph.person_id
		JOIN address_join aj ON p.id = aj.person_id
		JOIN address a ON aj.address_id = a.id
		WHERE p.id = ?;
	`

	row := db.DB.QueryRow(query, personID)
	err := row.Scan(&personInfo.Name, &personInfo.PhoneNumber, &personInfo.City, &personInfo.State, &personInfo.Street1, &personInfo.Street2, &personInfo.ZipCode)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": constants.PERSON_NOT_FOUND})
		} else {
			log.Fatal(err)
		}
		return
	}

	c.JSON(http.StatusOK, personInfo)
}

