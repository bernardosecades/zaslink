package main

import (
	"context"
	"fmt"
	"log"
	"time"

	sharesecretgrpc "github.com/bernardosecades/sharesecret/genproto"
	_ "github.com/bernardosecades/sharesecret/cmd"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:3333", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := sharesecretgrpc.NewSecretServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r1, err1 := client.CreateSecret(ctx, &sharesecretgrpc.CreateSecretRequest{Content: "this is a my secret"})
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
