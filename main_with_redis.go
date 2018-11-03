package main

import (
	"fmt"
	"time"
	"log"
	"os"

	"net/http"
	"math/rand"

	"github.com/go-redis/redis"
	
	"database/sql"
	_ "github.com/lib/pq"
)

var (
	dbHost     = envOrDefault("MYAPP_DATABASE_HOST", "localhost")
	dbPort     = envOrDefault("MYAPP_DATABASE_PORT", "5432")
	dbUser     = envOrDefault("MYAPP_DATABASE_USER", "root")
	dbPassword = envOrDefault("MYAPP_DATABASE_PASSWORD", "secret")
	dbName     = envOrDefault("MYAPP_DATABASE_NAME", "myapp")

	webHost = envOrDefault("MYAPP_WEB_HOST", "")
	webPort = envOrDefault("MYAPP_WEB_PORT", "8080")

	redisHost = envOrDefault("REDIS_HOST", "localhost")
	redisPort = envOrDefault("REDIS_PORT", "6379")

	db *sql.DB
	client *redis.Client
)

func envOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	fmt.Fprintln(w, "ID | Name")
	fmt.Fprintln(w, "---+--------")
	for rows.Next() {
		var (
			id   int
			name string
		)

		rows.Scan(&id, &name)

		fmt.Fprintf(w, "%2d | %s\n", id, name)
	}
}

func myCacheHandler(w http.ResponseWriter, r *http.Request) {
	n, err := client.Get("n").Result()
	if err != nil {
		panic(err)
		client.Set("n", rand.Intn(100), 5*time.Second)
	}

	fmt.Fprintf(w, "n = %d\n", n)
}

func NewClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:  redisHost + ":" + redisPort,
		Password: "", 
		DB:       0,  
	})

	if _, err := client.Ping().Result(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(client)
}

func main() {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	var err error
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	NewClient()

	http.HandleFunc("/", myHandler)
	http.HandleFunc("/cache", myCacheHandler)
	log.Print("Listening on " + webHost + ":" + webPort + "...")
	http.ListenAndServe(webHost+":"+webPort, nil)
}
