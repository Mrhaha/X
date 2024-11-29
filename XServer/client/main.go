package main

import (
	"XServer/serverproto/frame"
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := frame.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &frame.HelloRequest{Name: "World"}
	res, err := client.SayHello(ctx, req)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}

	log.Printf("Greeting: %s", res.GetMessage())
}
