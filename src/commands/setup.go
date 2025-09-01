package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

const githubWorkflow = `
name: Generate Changelog

on:
  push:
    tags:
      - 'v*'

jobs:
  changelog:
    runs-on: ubuntu-latest
    permissions:
      contents: write # Needed to commit back to the repo

    # IMPORTANT: Set these environment variables in your repository's secrets.
    # Settings > Secrets and variables > Actions > New repository secret
    #
    # - GCT_PROVIDER: The AI provider to use (e.g. OpenAI, Google AI Studio).
    # - GCT_MODEL: The specific model name (e.g. gpt-4o, gemini-1.5-flash).
    # - GCT_API_KEY: Your secret API key for the provider.
    #
    # You may also need to set provider-specific variables like GCT_AWS_REGION,
    # GCT_GCP_PROJECT_ID, etc., depending on your chosen provider.
    env:
      GCT_PROVIDER: ${{ secrets.GCT_PROVIDER }}
      GCT_MODEL: ${{ secrets.GCT_MODEL }}
      GCT_API_KEY: ${{ secrets.GCT_API_KEY }}

    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0 # Fetch all history for all tags

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Build GCT
        run: ./build/build.sh

      - name: Generate Changelog
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "41898282+github-actions[bot]@users.noreply.github.com"

          CURRENT_TAG=${{ github.ref_name }}
          echo "Current tag: $CURRENT_TAG"

          PREVIOUS_TAG=$(git tag --sort=-v:refname | sed -n '2p')
          if [ -z "$PREVIOUS_TAG" ]; then
            echo "No previous tag found, using first commit as base."
            PREVIOUS_TAG=$(git rev-list --max-parents=0 HEAD)
          fi
          echo "Previous tag/commit: $PREVIOUS_TAG"

          echo "## $CURRENT_TAG" > new_changelog.md
          ./gct ai log -c $PREVIOUS_TAG $CURRENT_TAG >> new_changelog.md
          echo "" >> new_changelog.md

          if [ -f Changelogs.md ]; then
            cat Changelogs.md >> new_changelog.md
          fi
          mv new_changelog.md Changelogs.md

      - name: Commit Changelog
        run: |
          git add Changelogs.md
          if git diff --staged --quiet; then
            echo "No changes to commit."
          else
            git commit -m "docs(changelog): Update Changelogs.md for ${{ github.ref_name }}"
            git push
          fi
`

const gitlabCI = `
workflow:
  rules:
    - if: $CI_COMMIT_TAG

stages:
  - changelog

generate_changelog:
  stage: changelog
  image: golang:1.22
  # IMPORTANT: Set these as CI/CD variables in your project's settings.
  # Settings > CI/CD > Variables > Add variable
  #
  # - GCT_PROVIDER: The AI provider to use (e.g. OpenAI, Google AI Studio).
  # - GCT_MODEL: The specific model name (e.g. gpt-4o, gemini-1.5-flash).
  # - GCT_API_KEY: Your secret API key. Make sure to set this as 'Masked'.
  #
  # You may also need to set provider-specific variables like GCT_AWS_REGION,
  # GCT_GCP_PROJECT_ID, etc., depending on your chosen provider.
  variables:
    GCT_PROVIDER: $GCT_PROVIDER
    GCT_MODEL: $GCT_MODEL
    GCT_API_KEY: $GCT_API_KEY

  before_script:
    - apt-get update -y && apt-get install -y git
    - git config --global user.name "${GITLAB_USER_NAME:-GitLab CI}"
    - git config --global user.email "${GITLAB_USER_EMAIL:-gitlab-ci@example.com}"
    - git remote set-url origin "https://gitlab-ci-token:${CI_JOB_TOKEN}@${CI_SERVER_HOST}/${CI_PROJECT_PATH}.git"
  script:
    - go build -o gct ./src
    - |
      CURRENT_TAG=$CI_COMMIT_TAG
      echo "Current tag: $CURRENT_TAG"

      PREVIOUS_TAG=$(git tag --sort=-v:refname | sed -n '2p')
      if [ -z "$PREVIOUS_TAG" ]; then
        echo "No previous tag found, using first commit as base."
        PREVIOUS_TAG=$(git rev-list --max-parents=0 HEAD)
      fi
      echo "Previous tag/commit: $PREVIOUS_TAG"

      echo "## $CURRENT_TAG" > new_changelog.md
      ./gct ai log -c $PREVIOUS_TAG $CURRENT_TAG >> new_changelog.md
      echo "" >> new_changelog.md

      if [ -f Changelogs.md ]; then
        cat Changelogs.md >> new_changelog.md
      fi
      mv new_changelog.md Changelogs.md

    - git add Changelogs.md
    - |
      if git diff --staged --quiet; then
        echo "No changes to commit."
      else
        git commit -m "docs(changelog): Update Changelogs.md for $CI_COMMIT_TAG"
        git push origin HEAD:${CI_DEFAULT_BRANCH:-main}
      fi
`

func SetupCommand() {
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	if len(os.Args) < 3 {
		fmt.Printf("%s Setup provider is required.\n", red("Error:"))
		fmt.Println("Usage: gct setup <github|gitlab>")
		return
	}
	provider := os.Args[2]

	var filePath, content string
	switch provider {
	case "github":
		dir := ".github/workflows"
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("%s Failed to create directory %s: %v\n", red("Error:"), dir, err)
			return
		}
		filePath = filepath.Join(dir, "changelog.yml")
		content = githubWorkflow
	case "gitlab":
		filePath = ".gitlab-ci.yml"
		content = gitlabCI
	default:
		fmt.Printf("%s Unsupported provider: %s. Use 'github' or 'gitlab'.\n", red("Error:"), provider)
		return
	}

	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("%s File '%s' already exists.\n", yellow("Warning:"), filePath)
		if !confirmPrompt("Do you want to overwrite it?") {
			fmt.Println("Setup cancelled.")
			return
		}
	}

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("%s Failed to write file %s: %v\n", red("Error:"), filePath, err)
		return
	}

	fmt.Printf("%s Successfully created '%s'.\n", green("✓"), filePath)
	fmt.Println("Please review the file and set the required environment variables (GCT_PROVIDER, GCT_MODEL, GCT_API_KEY) in your CI/CD provider's settings.")
}
