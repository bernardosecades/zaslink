package main

import (
	"fmt"
	"github.com/bernardosecades/sharesecret/handlers"
	"github.com/bernardosecades/sharesecret/repositories"
	"github.com/bernardosecades/sharesecret/services"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	commitHash string
)

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

	secretRepository := repositories.NewMySqlSecretRepository(dbName, dbUser, dbPass, dbHost, dbPort)
	secretService := services.NewSecretService(secretRepository, os.Getenv("SECRET_KEY"), os.Getenv("SECRET_PASSWORD"))
	secretHandler := handlers.NewSecretHandler(secretService)

	r := mux.NewRouter()

	r.HandleFunc("/secret/{id}", secretHandler.GetSecret).Methods("GET")
	r.HandleFunc("/secret", secretHandler.CreateSecret).Methods("POST")

	http.Handle("/", r)
	port := os.Getenv("SERVER_PORT")
	log.Print(fmt.Sprintf(":%s", port))
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
