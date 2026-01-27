package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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
		return nil
	} else {
		pathDirs := getPathDirs()
		for _, pathDir := range pathDirs {
			exists, absPath := getExecutableFromDir(command, pathDir)
			if exists {
				fmt.Printf("%s is %s\n", command, absPath)
				return nil
			}
		}
		fmt.Printf("%s: not found\n", command)
	}

	return nil
}

func getPathDirs() []string {
	return strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
}

func getExecutableFromDir(executableName, dirPath string) (exists bool, absPath string) {
	checkPath := filepath.Join(dirPath, executableName)
	fileInfo, err := os.Stat(checkPath)
	if err == nil && fileInfo.Mode().Perm()&os.FileMode(0o111) != 0 && !fileInfo.IsDir() {
		exists = true
		absPath = checkPath
	}
	return exists, absPath
}
