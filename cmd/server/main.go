package main

import (
	secretHandlers "github.com/bernardosecades/sharesecret/http"
	"github.com/bernardosecades/sharesecret/repository"
	"github.com/bernardosecades/sharesecret/service"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/joho/godotenv"

	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var commitHash string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Print("Not .env file found")
	}
}

func main() {

	fmt.Printf("Build Time: %s\n", time.Now().Format(time.RFC3339))
	fmt.Printf("Version: %s\n", commitHash)

	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")
	dbUser := os.Getenv("DB_USER")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	secretKey := os.Getenv("SECRET_KEY")
	secretPassword := os.Getenv("SECRET_PASSWORD")

	secretRepository := repository.NewMySQLSecretRepository(dbName, dbUser, dbPass, dbHost, dbPort)
	secretService := service.NewSecretService(secretRepository, secretKey, secretPassword)
	secretHandler := secretHandlers.NewSecretHandler(secretService)

	r := mux.NewRouter()

	r.HandleFunc("/secret/{id}", secretHandler.GetSecret).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/secret", secretHandler.CreateSecret).Methods(http.MethodPost, http.MethodOptions)

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"content-type"}),
		handlers.AllowedOrigins([]string{"*"}),
	)

	r.Use(cors)

	http.Handle("/", r)
	port := os.Getenv("SERVER_PORT")
	log.Print(fmt.Sprintf(":%s", port))
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
