package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"runtime"
)

var (
	resetColor  string
	grayColor   string
	lightBlue   string
	lightGreen  string
	lightRed    string
	lightYellow string
	lightCyan   string
)

type PromptConfig struct {
	ShowBranch           bool
	ShowModifiedFiles    bool
	ShowAheadBehind      bool
	ShowStashCount       bool
}

var defaultConfig = PromptConfig{
	ShowBranch:           true,
	ShowModifiedFiles:    false,
	ShowAheadBehind:      false,
	ShowStashCount:       false,
}

func init() {
	if runtime.GOOS == "windows" {
		resetColor = ""
		grayColor = ""
		lightBlue = ""
		lightGreen = ""
		lightRed = ""
		lightYellow = ""
		lightCyan = ""
	} else {
		resetColor = "\033[0m"
		grayColor = "\033[38;5;243m"
		lightBlue = "\033[38;5;117m"
		lightGreen = "\033[38;5;114m"
		lightRed = "\033[38;5;174m"
		lightYellow = "\033[38;5;186m"
		lightCyan = "\033[38;5;152m"
	}
}

func main() {
	config := loadConfig()

	for {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}

		prompt := buildPrompt(config)

		fmt.Printf("%s %s $ ", prompt, dir)

		var command string
		fmt.Scanln(&command)

		if command == "" {
			continue
		}

		args := strings.Fields(command)

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", command)
		} else {
			cmd = exec.Command(args[0], args[1:]...)
		}

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

func buildPrompt(config PromptConfig) string {
	var parts []string

	if config.ShowBranch {
		branch := getCurrentBranch()
		if branch != "" {
			status := getGitStatus()
			var branchColor string
			if status == "clean" {
				branchColor = lightGreen
			} else if status == "dirty" {
				branchColor = lightRed
			} else {
				branchColor = lightYellow
			}
			parts = append(parts, fmt.Sprintf("%s%s%s", branchColor, branch, grayColor))
		}
	}

	if config.ShowModifiedFiles {
		modifiedFiles := getModifiedFilesCount()
		parts = append(parts, fmt.Sprintf("%s%d%s", lightYellow, modifiedFiles, grayColor))
	}

	if config.ShowAheadBehind {
		aheadBehind := getAheadBehindStatus()
		if aheadBehind != "" {
			parts = append(parts, fmt.Sprintf("%s%s%s", lightBlue, aheadBehind, grayColor))
		}
	}

	if config.ShowStashCount {
		stashCount := getStashCount()
		parts = append(parts, fmt.Sprintf("%s%d%s", lightCyan, stashCount, grayColor))
	}

	if len(parts) > 0 {
		return fmt.Sprintf("%s(%s)%s", grayColor, strings.Join(parts, "|"), resetColor)
	}
	return ""
}

func loadConfig() PromptConfig {
	return defaultConfig
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

