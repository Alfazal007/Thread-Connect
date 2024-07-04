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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
	userRouter := routes.UserRouter(&apiCfg)
	r.Mount("/user", userRouter)
	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%v", portNumber),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("There was an error ", err)
	}
}
