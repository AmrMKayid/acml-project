package main

import (
	"database/sql"
	"encoding/json"
	"os"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Config struct {
    db []string
}

func main() {
    
    file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
	  fmt.Println("error:", err)
	}

	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		config.db.username, config.db.password, config.db.database)
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