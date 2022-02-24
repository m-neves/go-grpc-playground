package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/m-neves/go-grpc-playground/api/pb"
	"github.com/m-neves/go-grpc-playground/client/cmd"
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

		// Split incoming command into slice
		r := regexp.MustCompile(`\"[^\"]+\"|\S+`)
		m := r.FindAllString(c, -1)

		// Sanitize quotes on flags
		for i := range m {
			m[i] = strings.Trim(m[i], "\"")
		}

		if len(m) == 0 {
			log.Println("No command specified")
			continue
		}

		cmd.Exec(m, client)
	}
}
