package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		printPrompt()
		handleInput()
	}
}

func printPrompt() {
	fmt.Print("$ ")
}

func handleInput() {
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		os.Exit(1)
	}
	line = strings.TrimSpace(line)

	splitLine := strings.Split(line, " ")
	command := splitLine[0]
	args := splitLine[1:]

	if function, ok := builtins[command]; ok {
		function(args)
	} else {
		fmt.Printf("%s: command not found\n", command)
	}
}

type CommandHandler func(args []string) error

var builtins = map[string]CommandHandler{
	"exit": handleExit,
	"echo": handleEcho,
}

func handleEcho(args []string) error {
	fmt.Println(strings.Join(args, " "))
	return nil
}

func handleExit(args []string) error {
	os.Exit(0)
	return nil
}
