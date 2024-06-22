package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	var err error
	dsn := "postgresql://postgres:new+password@localhost:5432/cetec?sslmode=disable"
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	router.GET("/person/:person_id/info", getPersonInfo)
	router.POST("/person/create", createPerson)

	router.Run(":8080")
}

func getPersonInfo(c *gin.Context) {
	personID := c.Param("person_id")
	var person struct {
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
		City        string `json:"city"`
		State       string `json:"state"`
		Street1     string `json:"street1"`
		Street2     string `json:"street2"`
		ZipCode     string `json:"zip_code"`
	}

	query := `
        SELECT p.name, ph.number, a.city, a.state, a.street1, a.street2, a.zip_code
        FROM person p
        JOIN phone ph ON p.id = ph.person_id
        JOIN address_join aj ON p.id = aj.person_id
        JOIN address a ON aj.address_id = a.id
        WHERE p.id = $1
    `
	row := db.QueryRow(query, personID)
	err := row.Scan(&person.Name, &person.PhoneNumber, &person.City, &person.State, &person.Street1, &person.Street2, &person.ZipCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, person)
}

func createPerson(c *gin.Context) {
	var newPerson struct {
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
		City        string `json:"city"`
		State       string `json:"state"`
		Street1     string `json:"street1"`
		Street2     string `json:"street2"`
		ZipCode     string `json:"zip_code"`
	}

	if err := c.ShouldBindJSON(&newPerson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := tx.QueryRow("INSERT INTO person (name) VALUES ($1) RETURNING id", newPerson.Name)
	var personID int
	if err := result.Scan(&personID); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = tx.Exec("INSERT INTO phone (number, person_id) VALUES ($1, $2)", newPerson.PhoneNumber, personID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result = tx.QueryRow("INSERT INTO address (city, state, street1, street2, zip_code) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		newPerson.City, newPerson.State, newPerson.Street1, newPerson.Street2, newPerson.ZipCode)
	var addressID int
	if err := result.Scan(&addressID); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = tx.Exec("INSERT INTO address_join (person_id, address_id) VALUES ($1, $2)", personID, addressID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "person created"})
}
