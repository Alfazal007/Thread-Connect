package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"thread-connect/controllers"
	"thread-connect/internal/database"
	"thread-connect/routes"

	"github.com/cloudinary/cloudinary-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Error opening database connection", err)
	}
	apiCfg := controllers.ApiCfg{DB: database.New(conn)}
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	cloudApiKey := os.Getenv("CLOUDINARY_API_KEY")
	cloudApiSecret := os.Getenv("CLOUDINARY_API_SECRET")
	if cloudName == "" || cloudApiKey == "" || cloudApiSecret == "" {
		log.Fatal("There was an issue getting env variables of cloudinary")
	}

	cld, err := cloudinary.NewFromParams(cloudName, cloudApiKey, cloudApiSecret)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary, %v", err)
	}
	apiCfg.Cld = cld
	// REDIS CONFIGURATION
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPassword,
		DB:       0,
	})
	apiCfg.Rdb = rdb

	userRouter := routes.UserRouter(&apiCfg)
	r.Mount("/user", userRouter)
	tweetRouter := routes.TweetRouter(&apiCfg)
	r.Mount("/tweet", tweetRouter)
	followRouter := routes.FollowRouter(&apiCfg)
	r.Mount("/follow-user", followRouter)
	repostRouter := routes.RepostRouter(&apiCfg)
	r.Mount("/repost", repostRouter)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%v", portNumber),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("There was an error ", err)
	}
}