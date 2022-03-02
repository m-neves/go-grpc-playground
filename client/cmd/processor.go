package cmd

import (
	"context"
	"fmt"
	"log"

	"strings"

	"github.com/m-neves/go-grpc-playground/api/pb"
)

func Exec(cmd []string, c pb.GreetServiceClient) {
	// Interpreting command
	switch strings.ToLower(cmd[0]) {
	case "unary":
		if len(cmd) <= 1 {
			fmt.Println("Check README.md for usage details")
		} else {
			res, err := unary(cmd[1], c)

			if err != nil {
				log.Printf("Failed to invoke unary method: %s", err.Error())
				return
			}

			log.Println("Response:", res)
		}
	default:
		log.Println("Unknown command")
	}
}

func unary(message string, c pb.GreetServiceClient) (*pb.GreetResponse, error) {
	req := &pb.GreetRequest{Message: message}
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	return res, nil
}
