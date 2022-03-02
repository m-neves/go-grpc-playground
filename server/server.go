package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/m-neves/go-grpc-playground/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
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

func (s *server) GreetWithError(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error) {
	/*
		gRPC error docs:
		https://grpc.io/docs/guides/error/

		gRPC error codes:
		https://avi.im/grpc-errors/#go
	*/

	// Returns a error in case message is doerr
	msg := in.GetMessage()

	if msg == "doerr" {
		log.Println("returning a gRPC error response")
		return nil, status.Error(codes.Unauthenticated, "unauth gRPC error from server")
	}

	if msg == "err" {
		log.Println("returning a gRPC error response")
		return nil, errors.New("standard error")
	}

	log.Println("returning a successful gRPC call")
	return &pb.GreetResponse{Status: pb.ResponseStatus_SUCCESS}, nil
}

func (s *server) GreetWithTimeout(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error) {
	delay := 5

	msg := in.GetMessage()

	if msg == "ok" {
		delay = 1
	}

	for i := 0; i < delay; i++ {
		if ctx.Err() == context.Canceled {
			log.Println("Context cancelled")
			return nil, status.Error(codes.Canceled, "client cancelled the request")
		}

		time.Sleep(1 * time.Second)
	}

	return &pb.GreetResponse{Status: pb.ResponseStatus_SUCCESS}, nil
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
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	}

	log.Printf("Finished streaming %d messages", n)
	return nil
}

func (s *server) LongGreet(stream pb.GreetService_LongGreetServer) error {
	// Number of requests the server will accept
	const n = 10
	acceptedMessages := rand.Intn(n)

	// Processed messages is incremented after each successful iteration
	// As so, it starts as 1
	processedMessages := 1

	log.Printf("Processing %d requests", acceptedMessages)

	for {
		m, err := stream.Recv()

		if err == io.EOF {
			// If all messages get processed, return Success
			// Could break out of loop or return SendAndClose err
			log.Printf("All messages processed")
			return stream.SendAndClose(&pb.GreetResponse{
				Status: pb.ResponseStatus_SUCCESS,
			})
		}

		if err != nil {
			return err
		}

		var req pb.GreetRequest
		req.Message = m.GetMessage()

		log.Printf("Processed %v", &req)

		// End client streaming and send the client a response
		if processedMessages >= acceptedMessages {
			log.Printf("Processed %d messages", processedMessages)

			return stream.SendAndClose(&pb.GreetResponse{
				Status: pb.ResponseStatus_STREAM_END,
			})
		}

		processedMessages++
	}
}

func (s *server) GreetEveryone(stream pb.GreetService_GreetEveryoneServer) error {

	var responseStatus pb.ResponseStatus
	messagedToSend := 10 //rand.Intn(10)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Sending
	go func() {
		for i := 0; i <= messagedToSend; i++ {
			status := rand.Intn(len(pb.ResponseStatus_name))

			if status == 0 {
				status = 1
			}

			responseStatus = pb.ResponseStatus(status)
			stream.Send(&pb.GreetResponse{
				Status: responseStatus,
			})

			time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
		}

		wg.Done()
	}()

	// Receiving
	go func() {
		for {
			req, err := stream.Recv()

			if err == io.EOF {
				log.Println("Finished recieving")
				break
			}

			if err != nil {
				log.Println("Server receive error", err.Error())
				break
			}

			log.Printf("Recieved %v", req)
		}
		wg.Done()
	}()

	wg.Wait()

	fmt.Println("Finished")

	return nil
}

func main() {
	// Creates a TCP listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Fatalf("Failed to listen on port %d: %s", port, err.Error())
	}

	// TLS server
	certFile := "ssl/server.crt"
	keyFile := "ssl/server.pem"
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)

	if err != nil {
		log.Fatalf("Failed to start TLS Server: %v", err)
	}

	opts := grpc.Creds(creds)
	// End TLS

	// Creates a gRPC server
	s := grpc.NewServer(opts)

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
