package main

import (
	"XServer/serverproto/frame"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	frame.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *frame.HelloRequest) (*frame.HelloResponse, error) {
	log.Printf("Received: %v", req)
	return &frame.HelloResponse{Message: "Hello " + req.GetName()}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	frame.RegisterGreeterServer(grpcServer, &server{})

	log.Println("Server is listening on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
