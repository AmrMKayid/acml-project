package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PORT"), os.Getenv("PORT"), os.Getenv("PORT"))
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id   int
			name string
		)

		rows.Scan(&id, &name)

		log.Printf("%d: %s", id, name)
	}
}