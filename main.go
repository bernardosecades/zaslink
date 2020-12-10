package main

import (
	"fmt"
	"github.com/bernardosecades/sharesecret/controllers"
	"github.com/bernardosecades/sharesecret/repositories"
	"github.com/bernardosecades/sharesecret/services"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Print("Not .env file found")
	}
}

func main() {

	secretRepository := repositories.NewInMemorySecretRepository()
	secretService := services.NewSecretService(secretRepository, os.Getenv("SECRET_KEY"))
	secretController := controllers.NewSecretController(secretService)

	r := mux.NewRouter()

	r.HandleFunc("/secret/{id}", secretController.GetSecret).Methods("GET")
	r.HandleFunc("/secret", secretController.CreateSecret).Methods("POST")

	http.Handle("/", r)
	port := os.Getenv("SERVER_PORT")
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
