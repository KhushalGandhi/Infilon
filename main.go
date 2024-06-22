package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	var err error
	dsn := "username:password@tcp(127.0.0.1:3306)/cetec" // Change according to your local requirements
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	router.GET("/person/:person_id/info", getPersonInfo)

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
        WHERE p.id = ?
    `
	row := db.QueryRow(query, personID)
	err := row.Scan(&person.Name, &person.PhoneNumber, &person.City, &person.State, &person.Street1, &person.Street2, &person.ZipCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, person)
}
