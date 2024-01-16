package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Retrieve environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	tableName := os.Getenv("TABLE_NAME")

	if host == "" || port == "" || user == "" || password == "" || dbName == "" || tableName == "" {
		log.Fatal("DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, and TABLE_NAME environment variables are required.")
	}

	// Setup database connection
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the database")

	// Initialize database and tables
	initDatabase(db, dbName, tableName)

	// Initialize Gin router
	r := gin.Default()

	// Define a route that queries the database
	r.GET("/", func(c *gin.Context) {
		rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var result string
		for rows.Next() {
			var column1 string
			var column2 int
			// Replace column names and types with your actual table structure
			err := rows.Scan(&column1, &column2)
			if err != nil {
				log.Fatal(err)
			}
			result += fmt.Sprintf("Column1: %s, Column2: %d\n", column1, column2)
		}

		c.JSON(http.StatusOK, gin.H{"message": result})
	})

	// Define a route for inserting data into the table
	r.POST("/insert", func(c *gin.Context) {
		var requestData struct {
			Column1 string `json:"column1"`
			Column2 int    `json:"column2"`
		}

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Insert data into the table
		_, err := db.Exec(fmt.Sprintf("INSERT INTO %s (column1, column2) VALUES ($1, $2)", tableName), requestData.Column1, requestData.Column2)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Data inserted successfully"})
	})

	// Run the application
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

// Function to initialize database and tables
func initDatabase(db *sql.DB, dbName, tableName string) {
	// Create the database
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		log.Fatal(err)
	}

	// Switch to the created database
	_, err = db.Exec(fmt.Sprintf("USE %s", dbName))
	if err != nil {
		log.Fatal(err)
	}

	// Create the table if it doesn't exist
	_, err = db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			column1 VARCHAR(255),
			column2 INT
		)
	`, tableName))
	if err != nil {
		log.Fatal(err)
	}
}
