package commands

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type GitLabProvider struct{}

func NewGitLabProvider() (*GitLabProvider, error) {
	if _, err := exec.LookPath("glab"); err != nil {
		return nil, fmt.Errorf("'glab' (GitLab CLI) is not installed or not in your PATH")
	}
	return &GitLabProvider{}, nil
}

func (p *GitLabProvider) GetPRDetails(prNumber string) (*PRDetails, error) {
	cmd := exec.Command("glab", "mr", "view", prNumber, "--json", "title,body,author")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get MR details from glab CLI: %w. Is the MR number correct?", err)
	}

	var glabMR struct {
		Title  string `json:"title"`
		Body   string `json:"body"`
		Author struct {
			Username string `json:"username"`
		} `json:"author"`
	}

	if err := json.Unmarshal(output, &glabMR); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from glab CLI: %w", err)
	}

	diffCmd := exec.Command("glab", "mr", "diff", prNumber)
	diffOutput, err := diffCmd.Output()
	diffStr := ""
	if err != nil {
		diffStr = fmt.Sprintf("Could not fetch diff: %v", err)
	} else {
		diffStr = string(diffOutput)
	}

	return &PRDetails{
		Title:  glabMR.Title,
		Body:   glabMR.Body,
		Author: glabMR.Author.Username,
		Diff:   diffStr,
	}, nil
}

func (p *GitLabProvider) GetIssueDetails(issueNumber string) (*IssueDetails, error) {
	cmd := exec.Command("glab", "issue", "view", issueNumber, "--json", "title,body,author,labels")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get issue details from glab CLI: %w", err)
	}

	var glabIssue struct {
		Title  string `json:"title"`
		Body   string `json:"body"`
		Author struct {
			Username string `json:"username"`
		} `json:"author"`
		Labels []string `json:"labels"`
	}

	if err := json.Unmarshal(output, &glabIssue); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from glab CLI: %w", err)
	}

	return &IssueDetails{
		Title:  glabIssue.Title,
		Body:   glabIssue.Body,
		Author: glabIssue.Author.Username,
		Labels: glabIssue.Labels,
	}, nil
}
