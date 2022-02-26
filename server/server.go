package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"time"

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

func (s *server) GreetManyTimes(in *pb.GreetRequest, stream pb.GreetService_GreetManyTimesServer) error {
	const n = 10
	m := in.GetMessage()

	log.Printf("Streaming %d messages: %s", n, m)

	for i := 0; i <= n; i++ {
		res := &pb.GreetManyTimesResponse{
			Status:  pb.ResponseStatus_SUCCESS,
			Message: fmt.Sprintf("%s - %d", m, i),
		}

		stream.Send(res)
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	}

	log.Printf("Finished streaming %d messages", n)
	return nil
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
