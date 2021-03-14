package main

import (
	"log"
	"os"

	_ "github.com/bernardosecades/sharesecret/cmd"
	"github.com/bernardosecades/sharesecret/repository"
	"github.com/bernardosecades/sharesecret/service"
	"github.com/bernardosecades/sharesecret/server"
	"github.com/bernardosecades/sharesecret/server/grpc"
)

func main() {

	protocol := os.Getenv("SHARESECRET_SERVER_PROTOCOL")
	host     := os.Getenv("SHARESECRET_SERVER_HOST")
	port     := os.Getenv("SHARESECRET_SERVER_PORT")

	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")
	dbUser := os.Getenv("DB_USER")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	secretKey := os.Getenv("SECRET_KEY")
	secretPassword := os.Getenv("SECRET_PASSWORD")

	secretRepository := repository.NewMySQLSecretRepository(dbName, dbUser, dbPass, dbHost, dbPort)
	secretService := service.NewSecretService(secretRepository, secretKey, secretPassword)

	srvCfg := server.Config{Protocol: protocol, Host: host, Port: port}
	srv := grpc.NewServer(srvCfg, secretService)

	log.Printf("gRPC server running at %s://%s:%s ...\n", protocol, host, port)

	err := srv.Serve()

	if err != nil {
		log.Fatal("gRPC error: ", err)
	}
}
