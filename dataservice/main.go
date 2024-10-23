package main

import (
	"context"
	"net"
	"os"
	"log"

	pb "github.com/Amr2812/data-services/messages"
	"google.golang.org/grpc"
)

type MessageService struct {
	pb.UnimplementedMessagesServiceServer
}

func (s *MessageService) GetMessage(ctx context.Context, in *pb.MessageRequest) (*pb.MessageReply, error) {
	return &pb.MessageReply{}, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMessagesServiceServer(s, &MessageService{})
	log.Printf("Starting server on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		panic(err)
	}
}
