package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
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
	if len(splitLine) == 0 {
		return
	}
	command := splitLine[0]
	args := splitLine[1:]

	if function, ok := builtins[command]; ok {
		err := function(args)
		if err != nil {
			return
		}
	} else if strings.ContainsRune(command, '/') {
		fmt.Println("Error: direct paths not implemented")
	} else {
		exists, path := getExecutableFromPath(command)
		if exists {
			handleExecutables(command, path, args)
		} else {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}

type CommandHandler func(args []string) error

var builtins map[string]CommandHandler

func init() {
	builtins = map[string]CommandHandler{
		"exit": handleExit,
		"echo": handleEcho,
		"type": handleType,
		"pwd":  handlePwd,
		"cd":   handleCd,
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

func handlePwd(args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting pwd")
		return nil
	}

	fmt.Printf("%s\n", dir)
	return nil
}

func handleCd(args []string) error {
	if len(args) == 0 {
		return nil // home directory
	}

	dir := args[0]
	if err := os.Chdir(dir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("cd: %s: No such file or directory\n", dir)
		} else if errors.Is(err, syscall.ENOTDIR) {
			fmt.Printf("cd: %s: Not a directory\n", dir)
		} else if errors.Is(err, os.ErrPermission) {
			fmt.Printf("cd: %s: Permission Denied\n", dir)
		} else {
			fmt.Printf("cd: %s: %v\n", dir, err)
		}
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

func getExecutableFromPath(executableName string) (exists bool, absPath string) {
	pathDirs := getPathDirs()
	for _, pathDir := range pathDirs {
		exists, absPath := getExecutableFromDir(executableName, pathDir)
		if exists {
			return true, absPath
		}
	}
	return
}

func handleExecutables(commandName string, executablePath string, args []string) {
	executable := exec.Command(executablePath, args...)
	executable.Args[0] = commandName
	executable.Stdin = os.Stdin
	executable.Stdout = os.Stdout
	executable.Stderr = os.Stderr

	err := executable.Run()
	if err != nil {
		fmt.Println("Error running")
		return
	}
}
