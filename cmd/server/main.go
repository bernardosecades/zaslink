package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/bernardosecades/sharesecret/cmd"
	_ "github.com/bernardosecades/sharesecret/docs"
	secretHandlers "github.com/bernardosecades/sharesecret/http"
	"github.com/bernardosecades/sharesecret/repository"
	"github.com/bernardosecades/sharesecret/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title SHARE SECRET API
// @version 1.0
// @description API to create and read secrets one time

// @host
// @BasePath /
func main() {

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
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, fmt.Sprintf("%s/swagger/", os.Getenv("SERVER_URL")), http.StatusMovedPermanently)
	}).Methods(http.MethodGet, http.MethodOptions)

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

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
