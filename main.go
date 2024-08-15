package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	resetColor   = "\033[0m"
	boldColor    = "\033[1m"
	redColor     = "\033[31m"
	greenColor   = "\033[32m"
	yellowColor  = "\033[33m"
	magentaColor = "\033[35m"
	cyanColor    = "\033[36m"
)

func main() {
	for {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}

		branch := getCurrentBranch()
		status := getGitStatus()

		var branchColor string
		if status == "clean" {
			branchColor = greenColor
		} else if status == "dirty" {
			branchColor = redColor
		} else {
			branchColor = yellowColor
		}

		if branch != "" {
			fmt.Printf("%s(%s%s%s) %s $ ", cyanColor, branchColor, branch, resetColor, dir)
		} else {
			fmt.Printf("%s $ ", dir)
		}

		var command string
		fmt.Scanln(&command)

		args := strings.Split(command, " ")

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func getCurrentBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func getGitStatus() string {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	if len(output) == 0 {
		return "clean"
	}
	return "dirty"
}
