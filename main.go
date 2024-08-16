package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"runtime"
)

type Theme struct {
	ResetColor  string `json:"resetColor"`
	GrayColor   string `json:"grayColor"`
	LightBlue   string `json:"lightBlue"`
	LightGreen  string `json:"lightGreen"`
	LightRed    string `json:"lightRed"`
	LightYellow string `json:"lightYellow"`
	LightCyan   string `json:"lightCyan"`
}

type Config struct {
	PromptFormat string `json:"promptFormat"`
	Theme        Theme  `json:"theme"`
}

var (
	config Config
)

func init() {
	loadConfig()
	if runtime.GOOS == "windows" {
		config.Theme.ResetColor = ""
		config.Theme.GrayColor = ""
		config.Theme.LightBlue = ""
		config.Theme.LightGreen = ""
		config.Theme.LightRed = ""
		config.Theme.LightYellow = ""
		config.Theme.LightCyan = ""
	}
}

func loadConfig() {
	defaultConfig := Config{
		PromptFormat: "%s(%s%s%s|%s%d%s|%s%s%s|%s%d%s)%s %s $ ",
		Theme: Theme{
			ResetColor:  "\033[0m",
			GrayColor:   "\033[38;5;243m",
			LightBlue:   "\033[38;5;117m",
			LightGreen:  "\033[38;5;114m",
			LightRed:    "\033[38;5;174m",
			LightYellow: "\033[38;5;186m",
			LightCyan:   "\033[38;5;152m",
		},
	}

	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		config = defaultConfig
		return
	}

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		config = defaultConfig
	}
}

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
			branchColor = config.Theme.LightGreen
		} else if status == "dirty" {
			branchColor = config.Theme.LightRed
		} else {
			branchColor = config.Theme.LightYellow
		}

		if branch != "" {
			fmt.Printf(config.PromptFormat,
				config.Theme.GrayColor,
				branchColor, branch, config.Theme.GrayColor,
				config.Theme.LightYellow, modifiedFiles, config.Theme.GrayColor,
				config.Theme.LightBlue, aheadBehind, config.Theme.GrayColor,
				config.Theme.LightCyan, stashCount, config.Theme.GrayColor,
				config.Theme.ResetColor, dir)
		} else {
			fmt.Printf("%s $ ", dir)
		}

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

