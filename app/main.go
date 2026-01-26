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
	fmt.Printf("$ ")
}

func handleInput() {
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		os.Exit(1)
	}
	// line = strings.TrimSpace(line)

	splitLine := strings.Fields(line)
	command := splitLine[0]
	args := splitLine[1:]

	if function, ok := builtins[command]; ok {
		err := function(args)
		if err != nil {
			return
		}
	} else {
		fmt.Printf("%s: command not found\n", command)
	}
}

type CommandHandler func(args []string) error

var builtins map[string]CommandHandler

func init() {
	builtins = map[string]CommandHandler{
		"exit": handleExit,
		"echo": handleEcho,
		"type": handleType,
	}
}

func handleEcho(args []string) error {
	fmt.Println(strings.Join(args, " "))
	return nil
}

func handleExit(args []string) error {
	os.Exit(0)
	return nil
}

func handleType(args []string) error {
	command := args[0]
	if _, ok := builtins[command]; ok {
		fmt.Printf("%s is a shell builtin\n", args[0])
	} else {
		fmt.Printf("%s: not found\n", command)
	}
	return nil
}
