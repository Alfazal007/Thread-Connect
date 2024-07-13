package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"worker/controllers"
	"worker/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	portNumber := os.Getenv("PORT")
	if portNumber == "" {
		log.Fatal("Error getting port number from the env variables")
	}
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("Error getting database url from the env variables")
	}
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		log.Fatal("Error getting the redis url from the env variables")
	}
	redisPassword := os.Getenv("REDIS_PASSWORD") // TODO:: ADD ACTUAL PASSWORD AND UNCOMMENT THE LINE
	// if redisPassword == "" {
	// log.Fatal("Error getting the redis password from the env variables")
	// }
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Error opening database connection", err)
	}
	apiCfg := controllers.ApiCfg{DB: database.New(conn)}
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPassword,
		DB:       0,
	})
	apiCfg.Rdb = rdb
	for {
		data := apiCfg.Rdb.BRPop(context.Background(), 0, "worker")
		val, err := data.Result()
		if err != nil {
			fmt.Println("Error popping from Redis:", err)
			return
		}

		var repost controllers.Incoming
		err = json.Unmarshal([]byte(val[1]), &repost)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}
		if repost.Type == "repost" {
			apiCfg.ChangeRepostCount(repost)
		} else {
			println("Do something later")
		}
	}
}
