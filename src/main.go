package main

import (
	"fmt"
	"gct/src/commands"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	VerBranch = "Prod."
	VerStatus = "Release"
	VerNumber = "1.0.0"
	VerCommit = "dev"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "help" {
		commands.PrintUsage()
		return
	}

	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		commands.VersionCommand(VerBranch, VerStatus, VerNumber, VerCommit)
		return
	}

	if len(os.Args) >= 3 && os.Args[1] == "ai" && os.Args[2] == "commit" {
		additionalContext := ""
		if len(os.Args) > 3 {
			additionalContext = strings.Join(os.Args[3:], " ")
		}
		commands.AICommitCommand(additionalContext)
		return
	}

	if len(os.Args) >= 3 && os.Args[1] == "ai" && os.Args[2] == "diff" {
		commands.AIDiffCommand()
		return
	}

	if len(os.Args) >= 3 && os.Args[1] == "ai" && os.Args[2] == "log" {
		commands.AILogCommand()
		return
	}

	if len(os.Args) >= 3 && os.Args[1] == "ai" && os.Args[2] == "pr" {
		commands.AIPRCommand()
		return
	}

	if len(os.Args) >= 3 && os.Args[1] == "ai" && os.Args[2] == "issue" {
		commands.AIIssueCommand()
		return
	}

	if len(os.Args) >= 3 && os.Args[1] == "commit" && os.Args[2] == "edit" {
		commands.EditCommitCommand()
		return
	}

	if len(os.Args) >= 3 && os.Args[1] == "init" && os.Args[2] == "model" {
		commands.InitPresetCommand()
		return
	}

	command := os.Args[1]
	args := os.Args[2:]
	args2 := os.Args
	for _, arg := range args2 {
		if arg == "--no-cache" {
			commands.NoCache = true

			break
		}
	}

	switch command {
	case "init":
		if len(args) > 0 {
			if args[0] != "model" {
				fmt.Println(color.YellowString("Usage: gct init"))
				return
			}
		}
		commands.InitCommand()
	case "version":
		if len(args) > 0 {
			fmt.Println(color.YellowString("Usage: gct version (no arguments expected)"))
			return
		}
		commands.VersionCommand(VerBranch, VerStatus, VerNumber, VerCommit)
	case "about":
		if len(args) > 0 {
			fmt.Println(color.YellowString("Usage: gct about (no arguments expected)"))
			return
		}
		commands.AboutCommand()
	case "commit":
		if len(args) > 0 {
			fmt.Println(color.YellowString("Usage: gct commit (no arguments expected)"))
			return
		}
		commands.CommitCommand()
	case "ai":
		fmt.Printf("%s 'ai' command requires a subcommand.\n", color.RedString("Error:"))
		fmt.Println("Usage: gct ai [commit|diff]")
		return
	default:
		commands.NotFoundCommand()
	}
}
