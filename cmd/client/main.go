package main

import (
	"fmt"

	_ "github.com/bernardosecades/sharesecret/cmd"
	sharesecretgrpc "github.com/bernardosecades/sharesecret/build"

	"google.golang.org/grpc"


	"log"
	"time"
	"context"
)

// SEE: https://github.com/neocortical/mysvc/blob/main/grpc/client/client.go
func main() {
	conn, err := grpc.Dial("localhost:3333", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := sharesecretgrpc.NewSecretAppClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r1, err1 := client.CreateSecret(ctx, &sharesecretgrpc.CreateSecretRequest{Content: "this is a test"})
	if err1 != nil {
		log.Fatalf("CreateSecret: %v", err1)
	}
	fmt.Println(r1.String())

	r2, err2 := client.SeeSecret(ctx, &sharesecretgrpc.SeeSecretRequest{Id: r1.GetId()})

	if err2 != nil {
		log.Fatalf("SeeSecret: %v", err2)
	}
	fmt.Println(r2.String())
}

