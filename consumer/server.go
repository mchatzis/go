package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/mchatzis/go/producer/db/sqlc"
	producer_grpc "github.com/mchatzis/go/producer/grpc"
)

type server struct {
	producer_grpc.UnimplementedTaskServiceServer
	queries *sqlc.Queries
}

func (s *server) SendTask(ctx context.Context, task *producer_grpc.Task) (*emptypb.Empty, error) {
	// Here you can handle the incoming task, e.g., save it to the database
	log.Printf("Received task: %+v", task)
	return &emptypb.Empty{}, nil
}

func Listen() {
	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	producer_grpc.RegisterTaskServiceServer(s, &server{})

	log.Println("Server is running on port 8082...")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
