package service

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"

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
	case "bidi":
		biDiStream(c)
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

func biDiStream(c pb.GreetServiceClient) error {

	messages := []string{"Nothing", "but", "a", "southern", "soul", "deep", "inside", "of", "me", "Pride", "for"}
	messagesToSend := 9 //rand.Intn(10)

	client, err := c.GreetEveryone(context.Background())

	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Sending
	go func() {
		var req pb.GreetRequest
		for i := 0; i <= messagesToSend; i++ {
			req.Message = messages[i]
			client.Send(&req)

			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		client.CloseSend()
		wg.Done()
	}()

	// Receiving
	go func() {
		for {
			res, err := client.Recv()

			if err == io.EOF {
				log.Println("Finished receiving")
				break
			}

			if err != nil {
				// TODO implement
			}

			log.Printf("Received: %v", res)
		}

		wg.Done()
	}()

	wg.Wait()

	fmt.Println("Finished")
	return nil
}
