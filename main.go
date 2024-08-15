package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	for {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}

		branch := getCurrentBranch()

		if branch != "" {
			fmt.Printf("%s (%s) $ ", dir, branch)
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
