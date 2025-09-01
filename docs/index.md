---
title: GCT
description: An intelligent, AI-powered Git assistant.
---

[Repository](https://gitlab.com/Zillowe/Zillwen/Zusty/GCT)

## Features

- Conversational AI Commits: Generate a commit message and then "chat" with the AI to refine it until it's perfect.
- Multi-Provider Support: Works with over 10 AI providers, including OpenAI, Anthropic, Google (AI Studio & Vertex AI), Mistral, Amazon Bedrock, and any OpenAI-compatible endpoint.
- Git Hosting Integration: Explains pull/merge requests and proposes solutions for issues from GitHub, GitLab, and Forgejo.
- AI-Generated Changelogs: Automatically create user-facing changelogs from any set of git changes (gct ai log).
- AI-Powered Diff Analysis: Get a high-level explanation of any commit, branch, or staged changes (gct ai diff).
- Guided Setup: An interactive wizard (gct init) makes setup for any provider simple and fast.
- Automated Changelog Workflows: Generate and commit changelogs automatically on new tags with `gct setup`.
- Custom Guidelines: Enforce project-specific styles for commits and changelogs by providing your own guide files.

## Getting Started

The easiest way to get started is with the model preset wizard. It provides a curated list of popular, high-performance models and configures the provider for you. In your project directory, run:

```sh
gct init model
```

For a fully manual setup where you enter the provider and model name yourself, run `gct init`.

## Configuration

GCT can be configured via a local `gct.yaml` file, a global config file, or environment variables. The `init` command will help you create the local `gct.yaml` file, which holds the primary configuration for all AI commands.

**Security Note:** This file contains your API key. The `gct init` command will automatically add `gct.yaml` to your `.gitignore` file to prevent accidentally committing secrets.

### Supported AI Providers

GCT supports a wide range of AI providers:
`Google AI Studio`, `Google Vertex AI`, `OpenAI`, `OpenAI Compatible`, `Azure OpenAI`, `Anthropic`, `OpenRouter`, `DeepSeek`, `Mistral`, `Alibaba`, `Hugging Face`, `Amazon Bedrock`, and `xAI`.

### Supported Git Hosting Providers

The `ai pr` and `ai issue` commands integrate with the following platforms (via their respective CLIs):

- **GitHub** (via `gh`)
- **GitLab** (via `glab`)
- **Forgejo** and other Gitea forks (via `fj`)

### Configuration Fields

- `provider`: The AI service you want to use.
- `model`: The specific model/deployment name from your chosen provider.
- `api`: Your secret API key.
- `commits.guides` & `changelogs.guides`: Lists of local files containing formatting rules.
- `endpoint`: (Optional) The base URL, only for the `"OpenAI Compatible"` provider.
- `gcp_project_id`, `gcp_region`: (Optional) Required only for `"Google Vertex AI"`.
- `aws_region`, `aws_access_key_id`, `aws_secret_access_key`: (Optional) Required only for `"Amazon Bedrock"`.
- `azure_resource_name`: (Optional) Required only for `"Azure OpenAI"`.
- `cache.enabled`: (Optional) To store generated responses locally and reduce costs.

## Model Recommendations

Choosing a model can be tough. Here are some recommended starting points for GCT's use case:

| Recommendation           | Model ID                               | Best For...                                           |
| :----------------------- | :------------------------------------- | :---------------------------------------------------- |
| **Best Overall**         | `gpt-5-mini` (OpenAI)                  | Top-tier reasoning, speed, and instruction following. |
| **Best Balance**         | `claude-3-5-haiku-latest` (Anthropic)  | Excellent performance at a great price point.         |
| **Fastest & Best Value** | `gemini-2.5-flash` (Google)            | High-speed, low-cost tasks like `ai log` & chat.      |
| **Best Open Model**      | `openai/gpt-oss-120b` (via OpenRouter) | State-of-the-art open-source performance.             |

## Commands

GCT is a command-line tool. Here are the available commands, grouped by category:

### Core Commands

| Command                | Description                                                            |
| :--------------------- | :--------------------------------------------------------------------- |
| `gct init model`       | Starts a wizard with recommended models for easy setup.                |
| `gct init`             | Interactively creates a `gct.yaml` config file with manual input.      |
| `gct setup <provider>` | Creates a CI workflow (`github` or `gitlab`) for automated changelogs. |
| `gct version`          | Shows GCT version information.                                         |
| `gct help`             | Shows the detailed help message.                                       |

### Manual Git Commands

| Command           | Description                                             |
| :---------------- | :------------------------------------------------------ |
| `gct commit`      | Creates a new git commit using an interactive TUI form. |
| `gct commit edit` | Edits the previous commit's message using the same TUI. |

### AI Git Commands

| Command                   | Description                                                                  |
| :------------------------ | :--------------------------------------------------------------------------- |
| `gct ai commit [context]` | Generates and conversationally refines a commit message from staged changes. |
| `gct ai diff [args]`      | Asks AI to explain a set of code changes in a readable format.               |
| `gct ai log [args]`       | Generates a user-facing changelog entry from code changes.                   |
| `gct ai pr <number>`      | Summarizes a pull/merge request from GitHub, GitLab, or Forgejo.             |
| `gct ai issue <number>`   | Proposes a technical solution for an issue from GitHub, GitLab, or Forgejo.  |

## Installation

GCT must be built from source. To do so, you need to have [`go`](https://go.dev) installed.

```sh
# For Linux/macOS
./build/build-all.sh

# For Windows
./build/build-all.ps1
```

## Usage

GCT (Git Commit Tool) is a command-line interface designed to streamline your Git workflow with interactive forms and powerful AI integrations.

### Getting Started: The `init` Command

The first and most important step is to configure GCT for your project. This is required for all AI-powered commands. GCT offers two ways to create your configuration file:

#### Recommended: `gct init model`

This is the easiest way to get started. The `model` command launches a wizard that shows a curated list of high-performance, popular models. Simply choose a model, and GCT will configure the provider and model for you.

```sh
gct init model
```

#### Manual Setup: `gct init`

This command launches a fully manual, step-by-step setup wizard. It will ask you for every detail, including the provider name, the specific model ID, and any provider-specific information (like API endpoints or cloud project IDs).

```sh
gct init
```

Both commands will create a `gct.yaml` file in your project's root directory. For more details on the configuration options, please see the **[`gct.yaml` Configuration Reference](/docs/zds/gct/project-config)** page.

---

### Command Categories

GCT's commands are divided into three main categories:

- **Core Commands:** For managing the GCT application and configuration.
- **Manual Git Commands:** Interactive forms for common Git tasks that do not require AI.
- **AI-Powered Git Commands:** Commands that leverage AI to automate and enhance your workflow.

---

### Core Commands

- **`gct init model`**
  - Starts a guided setup wizard with a list of recommended models.
- **`gct init`**
  - Starts a fully manual setup wizard to create or overwrite the `gct.yaml` file.
- **`gct setup <github|gitlab>`**
  - Generates a CI/CD workflow file to automate changelog generation. When you push a new version tag (e.g. `v1.2.3`), the workflow will run, generate a changelog for the new version, and commit it to a `Changelogs.md` file in your repository.
  - **Usage:**
    - `gct setup github` (Creates `.github/workflows/changelog.yml`)
    - `gct setup gitlab` (Creates `.gitlab-ci.yml`)
- **`gct version`**
  - Shows the currently installed GCT version and build details.
- **`gct about`**
  - Displays information about the GCT project.
- **`gct help`**
  - Shows the detailed help message listing all available commands.

---

### Manual Git Commands

- **`gct commit`**
  - Launches a full-screen interactive form to guide you through creating a well-structured commit message. It prompts for a `Type`, `Subject`, and `Body`, then asks for confirmation before committing.

- **`gct commit edit`**
  - Allows you to easily amend the _most recent_ commit's message. It fetches the last message and pre-populates the same interactive form from `gct commit`.

---

### AI-Powered Git Commands

- **`gct ai commit [context]`**
  - This is the flagship feature. It automatically generates a commit message from your staged changes and then allows you to perfect it conversationally.
  - **Workflow:**
    1.  The AI generates a commit message based on your staged code and any guidelines in your `gct.yaml`.
    2.  The message is displayed for your review.
    3.  You are prompted with options:
        - **[c] to chat/change:** Provide a follow-up instruction (e.g. "add a co-author," "make the subject shorter") and the AI will revise the message.
        - **[e] to edit:** Open the generated message in the manual TUI editor for full control.
        - **[Enter] to commit:** Accept the message and commit it directly.
        - **[q] to quit:** Cancel the operation.
  - **Providing Context:** You can give the AI extra information by passing it as an argument:
    ```sh
    gct ai commit "This change was co-authored by Jane Doe and fixes issue #123."
    ```

- **`gct ai diff [arguments]`**
  - Asks an AI to act as an expert code reviewer, providing a high-level explanation of code changes. The output is displayed in a clean, scrollable TUI.
  - **Usage Examples:**
    - `gct ai diff` (Explains unstaged changes)
    - `gct ai diff --staged` (Explains staged changes)
    - `gct ai diff <commit-hash>` (Explains a specific commit)
    - `gct ai diff <branch-name>` (Explains changes relative to another branch)

- **`gct ai log [arguments]`**
  - Generates a user-facing changelog entry from a set of code changes. It uses the same arguments as `gct ai diff` but provides output formatted for release notes.
  - **Usage Examples:**
    - `gct ai log` (Creates a changelog for unstaged changes)
    - `gct ai log --staged` (Creates a changelog for staged changes)
    - `gct ai log <commit-hash>` (Creates a changelog for a specific commit)
    - `gct ai log <branch-name>` (Creates a changelog for changes on a branch)
    - `gct ai log v1.0.0 v1.1.0` (Creates a changelog for changes between two tags)
  - **Non-Interactive Output:**
    - For use in CI/CD pipelines or scripts, add the `-c` flag to print the raw markdown output directly to the console without the interactive viewer.
    - `gct ai log -c v1.0.0 v1.1.0`

- **`gct ai pr <number>`**
  - Summarizes a pull request or merge request from a supported git hosting provider (GitHub, GitLab, Forgejo). It provides a high-level overview of the changes, the purpose, and the solution.
  - **Usage:**
    - `gct ai pr 123`

- **`gct ai issue <number>`**
  - Proposes a technical implementation plan for an issue from a supported git hosting provider. It outlines the "why", the "how", and the "solution".
  - **Usage:**
    - `gct ai issue 456`

---

### Global Flags

- **`gct -v`**, **`gct --version`**
  - A global alternative to the `gct version` command.

## FAQ

Frequently Asked Questions

<Accordions type="single">
  <Accordion title="What is GCT?">
    GCT is an advanced command-line tool that enhances your Git workflow. It started as an interactive tool for writing well-structured commit messages but has evolved into a powerful AI-assisted development partner. It can automatically generate and conversationally refine commit messages, explain complex code changes, create changelogs, and more—all from your terminal.
  </Accordion>
</Accordions>
<br />
<Accordions type="single">
  <Accordion title="How do I use GCT for free?">
    You can use GCT for free by setting the provider as OpenRouter and use a free model from OpenRouter, or to use Google AI Studio and use a model like Gemini 2.5 Flash
  </Accordion>
</Accordions>
<br />
<Accordions type="single">
  <Accordion title="What's the difference between 'gct init' and 'gct init model'?">
    Both commands create your `gct.yaml` configuration file, but they offer different experiences:
    - **`gct init model` (Recommended):** This is the easiest way to start. It presents a curated list of popular, high-performance models. You just pick one from the list, and GCT will configure the provider and model name for you.
    - **`gct init` (Manual):** This gives you full control. It prompts you to manually enter the provider name, model name, and any other provider-specific details. You should use this when you want to use a model not in the model list or need to configure a complex provider from scratch.
  </Accordion>
  </Accordions>
<br />

<Accordions type="single">
  <Accordion title="Do I need to use the AI features?">
    Not at all! GCT is designed to be useful for everyone:
    1. **As a Manual Git Helper:** Commands like `gct commit` and `gct commit edit` provide an interactive form (TUI) to help you write perfectly formatted commit messages without any AI.
    2. **As an AI Assistant:** Commands like `gct ai commit`, `gct ai diff`, and `gct ai log` use AI to automate tasks. These features are entirely optional and require a `gct.yaml` config file.
  </Accordion>
</Accordions>
<br />
<Accordions type="single">
  <Accordion title="What is the 'conversational' part of 'gct ai commit'?">
   After the AI generates the first draft of a commit message, GCT doesn't just ask you to accept or reject it. It gives you a "chat" option. You can provide follow-up instructions like "add a co-author," "make the subject line shorter," or "explain the performance impact in the body," and the AI will revise the message based on your feedback. You can do this as many times as you need until the message is perfect.
  </Accordion>
</Accordions>
<br />
<Accordions type="single">
  <Accordion title="Which AI providers are supported?">
    GCT supports a vast and growing ecosystem of AI services. You can connect to almost any major platform:
    - **Major Platforms:** OpenAI, Anthropic (Claude), Google AI Studio (Gemini), Google Vertex AI, Azure OpenAI, Amazon Bedrock.
    - **Open Model Providers:** Mistral, DeepSeek, Alibaba (Qwen), xAI (Grok).
    - **Aggregators & Endpoints:** OpenRouter, Hugging Face, and any other "OpenAI Compatible" endpoint.
  </Accordion>
</Accordions>
<br />
<Accordions type="single">
  <Accordion title="How does the AI know how to format my commits and changelogs?">
    The `gct.yaml` file has two sections for providing style guides: `commits` and `changelogs`. You can list paths to local Markdown (`.md`) or text (`.txt`) files in each section.
    - When you run `gct ai commit`, it uses the `commits.guides`.
    - When you run `gct ai log`, it uses the `changelogs.guides`.
      This ensures the AI's output is always tailored to your specific project conventions.
    </Accordion>
</Accordions>
<br />
<Accordions type="single">
  <Accordion title="How does authentication work for providers like AWS Bedrock or Google
Vertex AI?">
    Instead of a single API key, these enterprise platforms use more complex authentication, which the `gct init` wizard will guide you through:
    - For **Amazon Bedrock**, it will prompt for your AWS Access Key ID, Secret AccessKey, and AWS Region.
    - For **Google Vertex AI**, it will prompt for your GCP Project ID and GCP Region. It then uses your machine's existing Google Cloud authentication (usually configured by running `gcloud auth application-default login`).
    - For all other providers, it will ask for a standard API key.
  </Accordion>
</Accordions>
<br />
<Accordions type="single">
  <Accordion title="How is my API key stored? Is it secure?">
   Your API keys and credentials are stored in plain text in the `gct.yaml` file within your project's directory.
  **This is why it is critical that you add `gct.yaml` to your project's `.gitignore` file.** The `gct init` command does this for you automatically to help prevent your secret credentials from ever being committed to your repository.
  </Accordion>
</Accordions>
<br />
<Accordions type="single">
  <Accordion title="What happens if a Git command executed by GCT fails?">
    If the underlying Git command fails (e.g. `git commit` with no staged changes, or a failing pre-commit hook), GCT will catch the error. It will display a message indicating the failure and will print the exact command that was attempted, which is very helpful for debugging the issue.
  </Accordion>
</Accordions>
