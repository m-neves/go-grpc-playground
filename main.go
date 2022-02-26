package main

import (
	"flag"
	"fmt"
	"strings"
)

type Command struct {
	command string
	message string
}

func main() {
	test1()
}

func test1() {
	args := []string{"-arg", "-flag", "value"}
	flags := flag.NewFlagSet("flags", flag.ExitOnError)

	flagBool := flags.Bool("arg", false, "")
	flagValue := flags.String("flag", "defaultValue", "")

	flags.Parse(args)

	fmt.Println(*flagValue)
	fmt.Println(*flagBool)
}

func test2() {
	command := "unary -flag apu nahasa"

	args := strings.Split(command, " ") //[]string{"-flag", "value nahasa", "arg"}
	args = []string{args[0], args[1], strings.Join(args[2:], " ")}

	flags := flag.NewFlagSet("flags", flag.ExitOnError)

	unary := flags.Bool("unary", false, "")
	flagValue := flags.String("flag", "defaultValue", "")

	flags.Parse(args)

	fmt.Println(*flagValue)

	fmt.Println(*unary)
}
