package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"math/rand"

	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"

	"github.com/go-redis/cache"
	
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

	codec
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
	// random n
	err := client.Set("key", "value", rand.Intn(100)).Err()
	if err != nil {
		panic(err)
	}

	n, err := client.Get("key").Result()
	if err != nil {
		panic(err)
	}

	codec.Set(&cache.Item{
		Key:        "key",
		Object:     n,
		Expiration: 5*time.Second,
	})


	fmt.Fprintln(w, "n = " + n)
}

func NewClient() {
	client := redis.NewClient(&redis.Options{
		Addr:  redisHost + ":" + redisPort,
		Password: "", 
		DB:       0,  
	})

	if err = client.Ping().Result(); err != nil {
		log.Fatal(err)
	}

	codec := &cache.Codec{
		Redis: client,

		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
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
