package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	resetColor  = "\033[0m"
	grayColor   = "\033[38;5;243m"
	lightBlue   = "\033[38;5;117m"
	lightGreen  = "\033[38;5;114m"
	lightRed    = "\033[38;5;174m"
	lightYellow = "\033[38;5;186m"
	lightCyan   = "\033[38;5;152m"
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
		modifiedFiles := getModifiedFilesCount()
		aheadBehind := getAheadBehindStatus()
		stashCount := getStashCount()

		var branchColor string
		if status == "clean" {
			branchColor = lightGreen
		} else if status == "dirty" {
			branchColor = lightRed
		} else {
			branchColor = lightYellow
		}

		if branch != "" {
			fmt.Printf("%s(%s%s%s|%s%d%s|%s%s%s|%s%d%s)%s %s $ ",
				grayColor,
				branchColor, branch, grayColor,
				lightYellow, modifiedFiles, grayColor,
				lightBlue, aheadBehind, grayColor,
				lightCyan, stashCount, grayColor,
				resetColor, dir)
		} else {
			fmt.Printf("%s $ ", dir)
		}

		var command string
		fmt.Scanln(&command)

		if command == "" {
			continue
		}

		args := strings.Split(command, " ")

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			if _, ok := err.(*exec.ExitError); !ok {
				fmt.Println("Command not found:", args[0])
			}
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

func getModifiedFilesCount() int {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}
	return len(strings.Split(strings.TrimSpace(string(output)), "\n"))
}

func getAheadBehindStatus() string {
	cmd := exec.Command("git", "rev-list", "--left-right", "--count", "HEAD...@{u}")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	counts := strings.Split(strings.TrimSpace(string(output)), "\t")
	if len(counts) != 2 {
		return ""
	}
	ahead, _ := strconv.Atoi(counts[0])
	behind, _ := strconv.Atoi(counts[1])
	return fmt.Sprintf("↑%d↓%d", ahead, behind)
}

func getStashCount() int {
	cmd := exec.Command("git", "stash", "list")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}
	return len(strings.Split(strings.TrimSpace(string(output)), "\n"))
}
