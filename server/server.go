package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/m-neves/go-grpc-playground/api/pb"
	"google.golang.org/grpc"
)

const port = 5000

type server struct {
	pb.UnimplementedGreetServiceServer
}

func (s *server) Greet(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error) {
	log.Printf("Greet invoked with message %q", in.GetMessage())
	in.GetMessage()

	res := &pb.GreetResponse{Status: pb.ResponseStatus_SUCCESS}

	return res, nil
}

func main() {
	// Creates a TCP listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Fatalf("Failed to listen on port %d: %s", port, err.Error())
	}

	// Creates a gRPC server
	s := grpc.NewServer()

	// Register the generated protobuf service
	pb.RegisterGreetServiceServer(s, &server{})

	log.Println("Serving on port", port)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve server on port %d: %s", port, err.Error())
		}
	}()

	<-sig

	log.Printf("Gracefully shutting down server")
	s.GracefulStop()
}
