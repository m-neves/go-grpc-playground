package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/m-neves/go-grpc-playground/api/pb"
	"github.com/m-neves/go-grpc-playground/client/service"

	"google.golang.org/grpc"
)

const port = 5000

func main() {
	conn, err := grpc.Dial(fmt.Sprintf(":%d", port), grpc.WithInsecure())

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
