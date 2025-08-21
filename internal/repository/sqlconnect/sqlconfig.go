package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func ConnectDb(dbname string) (*sql.DB, error) {

	fmt.Println("Connecting to database:", dbname)
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file", err)
		return nil, err
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", user, password, dbHost, dbPort, dbname)

	fmt.Println(connectionString)
	// Open database connection
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify the connection is actually working.
	err = db.Ping()
	if err != nil {
		db.Close() // Close the connection if ping fails
		return nil, err
	}

	fmt.Println("Connected to database:", dbname)
	return db, nil
}
