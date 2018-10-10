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
    Db struct {
        Username     string `json:"username"`
        Password string `json:"password"`
        Database string `json:"database"`
    } `json:"db"`
}

func main() {
	config := LoadConfiguration("config.json")
	
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		config.Db.Username, config.Db.Password, config.Db.Database)
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

func LoadConfiguration(file string) Config {
    var config Config
    configFile, err := os.Open(file)
    defer configFile.Close()
    if err != nil {
        fmt.Println(err.Error())
    }
    jsonParser := json.NewDecoder(configFile)
    jsonParser.Decode(&config)
    return config
}