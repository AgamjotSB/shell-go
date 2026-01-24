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
	args := strings.Join(splitLine[1:], " ")

	switch command {
	case "exit":
		os.Exit(0)
	case "echo":
		fmt.Println(args)

	default:
		fmt.Printf("%s: command not found\n", command)
	}
}
