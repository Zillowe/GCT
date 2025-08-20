---
title: Project Configuration
description: Reference for the gct.yaml configuration file.
---

# `gct.yaml` Configuration Reference

The `gct.yaml` file is the heart of GCT's AI capabilities. It tells the tool which AI provider to use, which model to call, and how to authenticate. It also allows you to provide custom guidelines to tailor the AI's output to your project's specific needs.

The `gct init` and `gct init model` commands will generate this file for you automatically.

## Top-Level Fields

| Field      | Type     | Description                                                                                                                                                  |
| :--------- | :------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `name`     | `string` | A friendly name for your project.                                                                                                                            |
| `provider` | `string` | The AI provider to use (e.g. `OpenAI`, `Anthropic`, `Google AI Studio`). Must match one of the [Supported Providers](/docs/zds/gct/#supported-ai-providers). |
| `model`    | `string` | The specific model or deployment ID for the chosen provider (e.g. `gpt-4o`, `claude-3-5-haiku-latest`).                                                      |
| `api`      | `string` | Your secret API key for the chosen provider. **This is a secret and should not be committed.**                                                               |
| `cache`    | `object` | (Optional) Settings for caching AI responses to reduce costs and latency.                                                                                    |

## Provider-Specific Fields

Some providers require additional information beyond an API key. These fields should only be included when using the specified provider.

| Field                   | Provider            | Description                                                             |
| :---------------------- | :------------------ | :---------------------------------------------------------------------- |
| `endpoint`              | `OpenAI Compatible` | The base URL of the API endpoint (e.g. `https://api.example.com/v1`).   |
| `gcp_project_id`        | `Google Vertex AI`  | Your Google Cloud Platform Project ID.                                  |
| `gcp_region`            | `Google Vertex AI`  | The GCP region for your Vertex AI model (e.g. `us-central1`).           |
| `aws_region`            | `Amazon Bedrock`    | The AWS region where your Bedrock models are hosted (e.g. `us-east-1`). |
| `aws_access_key_id`     | `Amazon Bedrock`    | Your AWS Access Key ID for authentication.                              |
| `aws_secret_access_key` | `Amazon Bedrock`    | Your AWS Secret Access Key for authentication.                          |
| `azure_resource_name`   | `Azure OpenAI`      | The name of your Azure OpenAI resource.                                 |

## Custom Guidelines

You can provide custom instructions to the AI to ensure its output matches your project's conventions. This is useful for enforcing commit message formats (like Conventional Commits) or changelog styles.

| Field               | Type    | Description                                                                                         |
| :------------------ | :------ | :-------------------------------------------------------------------------------------------------- |
| `commits.guides`    | `array` | A list of paths to local `.md` or `.txt` files that will be used as guidelines for `gct ai commit`. |
| `changelogs.guides` | `array` | A list of paths to local `.md` or `.txt` files that will be used as guidelines for `gct ai log`.    |

### Example `gct.yaml`

```yaml
name: My Awesome Project
provider: OpenAI
model: gpt-4o
api: sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
commits:
  guides:
    - ./docs/guides/commits.md
changelogs:
  guides: []
cache:
  enabled: true
```
