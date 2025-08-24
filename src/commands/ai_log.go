package commands

import (
	"fmt"
	"gct/src/config"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

const aiLogPromptTemplateWithGuidelines = `
You are a release manager writing a changelog. Based on the following git diff and guidelines, generate a concise and user-friendly changelog entry.

Here are the guidelines to follow:
--- GUIDELINES START ---
%s
--- GUIDELINES END ---

Follow these rules:
1.  **Structure:** Use Markdown with headings for different categories (e.g. ### ✨ Features, ### 🌟 Enhancements, ### 🐛 Bug Fixes).
2.  **Clarity:** Write in the present tense (e.g. "Add feature" not "Added feature").
3.  **Focus:** Emphasize user-facing changes. Ignore minor code-quality improvements or refactoring unless they have a direct impact.
4.  **Conciseness:** Use bullet points for individual changes.

--- GIT DIFF START ---
%s
--- GIT DIFF END ---

Provide only the Markdown for the changelog entry below:
`

const aiLogPromptTemplate = `
You are a release manager writing a changelog. Based on the following git diff, generate a concise and user-friendly changelog entry.

Follow these rules:
1.  **Structure:** Use Markdown with headings for different categories (e.g. ### ✨ Features, ### 🐛 Bug Fixes, ### 🚀 Performance).
2.  **Clarity:** Write in the present tense (e.g. "Add feature" not "Added feature").
3.  **Focus:** Emphasize user-facing changes. Ignore minor code-quality improvements or refactoring unless they have a direct impact.
4.  **Conciseness:** Use bullet points for individual changes.

--- GIT DIFF START ---
%s
--- GIT DIFF END ---

Provide only the Markdown for the changelog entry below:
`

func AILogCommand() {
	cyan := color.New(color.FgCyan).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	var diffCmd *exec.Cmd
	var description string
	isCI := false
	args := os.Args[3:]

	if len(args) > 0 && (args[0] == "-c" || args[0] == "--ci") {
		isCI = true
		args = args[1:]
	}

	switch {
	case len(args) == 0:
		description = "unstaged changes"
		diffCmd = exec.Command("git", "diff")
	case len(args) == 1 && args[0] == "--staged":
		description = "staged changes"
		diffCmd = exec.Command("git", "diff", "--staged")
	case len(args) == 1:
		ref := args[0]
		description = fmt.Sprintf("changes from '%s'", ref)
		diffCmd = exec.Command("git", "diff", ref)
	case len(args) == 2:
		startTag, endTag := args[0], args[1]
		description = fmt.Sprintf("changes between '%s' and '%s'", startTag, endTag)
		diffCmd = exec.Command("git", "diff", fmt.Sprintf("%s..%s", startTag, endTag))
	default:
		fmt.Printf("%s Invalid arguments for 'ai log'.\n", red("Error:"))
		fmt.Println("Usage: gct ai log [-c] [--staged | <commit|branch> | <start_tag> <end_tag>]")
		return
	}

	if !isCI {
		fmt.Printf("%s Generating changelog for %s...\n", cyan("🔍"), description)
	}

	diffOutput, err := diffCmd.Output()
	if err != nil {
		fmt.Printf("%s Failed to get git diff. Is the reference valid?\n", red("Error:"))
		return
	}
	if len(diffOutput) == 0 {
		if !isCI {
			fmt.Printf("%s No changes found to generate a changelog for %s.\n", green("✓"), description)
		}
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("%s %v\n", color.RedString("Error:"), err)
		return
	}

	guidelines, _ := readGuidelines(cfg.Changelogs.Paths)

	var prompt string
	if guidelines != "" {
		if !isCI {
			fmt.Println(cyan("📚 Reading changelog guidelines..."))
		}
		prompt = fmt.Sprintf(aiLogPromptTemplateWithGuidelines, guidelines, string(diffOutput))
	} else {
		prompt = fmt.Sprintf(aiLogPromptTemplate, string(diffOutput))
	}
	aiResponse, err := runAITask(prompt, isCI)
	if err != nil {
		fmt.Printf("%s %v\n", red("Error:"), err)
		return
	}

	cleanMsg := strings.TrimSpace(aiResponse)

	if isCI {
		fmt.Println(cleanMsg)
	} else {
		viewerModel := NewAITextViewerModel("🤖 AI Generated Changelog", cleanMsg)
		p := tea.NewProgram(viewerModel, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Printf("%s Error displaying AI response: %v\n", red("Error:"), err)
		}
	}
}
