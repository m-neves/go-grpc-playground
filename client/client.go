package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/m-neves/go-grpc-playground/api/pb"
	"github.com/m-neves/go-grpc-playground/client/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const port = 5000

func main() {
	// TLS
	certFile := "ssl/ca.crt"
	creds, err := credentials.NewClientTLSFromFile(certFile, "")

	if err != nil {
		log.Fatalf("Failed to start TLS client: %v", err)
	}

	opts := grpc.WithTransportCredentials(creds)
	// End TLS

	conn, err := grpc.Dial(fmt.Sprintf(":%d", port), opts)

	if err != nil {
		log.Fatalf("Failed to dial at port %d: %s", port, err.Error())
	}

	defer conn.Close()

	c := pb.NewGreetServiceClient(conn)

	go readConsole(c)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig
	fmt.Println("Bye")
}

func readConsole(client pb.GreetServiceClient) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		c, _ := reader.ReadString('\n')

		// Convert CRLF to LF
		c = strings.Replace(c, "\n", "", -1)

		service.Exec(c, client)
	}
}

func Unary(message string, c pb.GreetServiceClient) (*pb.GreetResponse, error) {
	req := &pb.GreetRequest{Message: message}
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	return res, nil
}
