package service

import (
	"context"
	"flag"
	"io"
	"log"

	"strings"

	"github.com/m-neves/go-grpc-playground/api/pb"
)

type Command struct {
	Command string
	Message string
}

func Exec(cmd string, c pb.GreetServiceClient) {
	var command Command
	args := strings.Split(cmd, " ")

	// Concatenates the remainder values as -msg parameter
	args = []string{args[0], strings.Join(args[1:], " ")}

	flags := flag.NewFlagSet("command", flag.ContinueOnError)

	flags.StringVar(&command.Command, "cmd", "", "")
	flags.StringVar(&command.Message, "msg", "", "")

	err := flags.Parse(args)

	if err != nil {
		log.Printf("Failed to parse command: %s", cmd)
	}

	switch command.Command {
	case "unary":
		unary(command.Message, c)
	case "sstream":
		serverStream(command.Message, c)
	case "cstream":
		clientStream(c)
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

func serverStream(message string, c pb.GreetServiceClient) error {
	req := &pb.GreetRequest{Message: message}
	stream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			log.Println("End of stream")
			break
		}

		if err != nil {
			return err
		}

		log.Printf("Received message: %s", msg)
	}

	return nil
}

func clientStream(c pb.GreetServiceClient) error {
	messages := []string{"Devilish Trio", "Drake", "Linkin Park", "Bones"}

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		return err
	}

	var req pb.GreetRequest
	for _, v := range messages {
		req.Message = v
		stream.Send(&req)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Printf("Received from server: %v", res)

	return nil
}
