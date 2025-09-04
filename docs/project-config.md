---
title: Project Configuration
description: Reference for the gct.yaml configuration file and environment variables.
---

GCT offers a flexible configuration system that allows you to set up your AI provider and other settings through files and environment variables.

## Configuration Hierarchy

GCT loads configuration from the following sources in order, with later sources overriding earlier ones:

1.  **Global Config File**: A `config.yaml` file located in your user config directory. This is useful for setting a default provider and API key that you use across all your projects.
    - **Linux/macOS**: `~/.config/gct/config.yaml` (or `~/.gct/config.yaml` for legacy)
    - **Windows**: `%APPDATA%\gct\config.yaml`
2.  **Local Config File**: A `gct.yaml` file in your project's directory (or any parent directory). This allows you to have project-specific settings, like custom commit guidelines.
3.  **Environment Variables**: Any variable prefixed with `GCT_`. These are perfect for CI/CD environments or for temporarily overriding a setting without modifying a file. You can also place these in a `.env` file in your project directory for them to be loaded automatically.

---

## `gct.yaml` File Reference

The `gct.yaml` file is the heart of GCT's AI capabilities. The `gct init` and `gct init model` commands will generate this file for you automatically.

### Top-Level Fields

| Field      | Type     | Required | Description                                                                                                                                                  |
| :--------- | :------- | :------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `name`     | `string` | No       | A friendly name for your project.                                                                                                                            |
| `provider` | `string` | **Yes**  | The AI provider to use (e.g. `OpenAI`, `Anthropic`, `Google AI Studio`). Must match one of the [Supported Providers](/docs/zds/gct/#supported-ai-providers). |
| `model`    | `string` | **Yes**  | The specific model or deployment ID for the chosen provider (e.g. `gpt-4o`, `claude-3-5-haiku-latest`).                                                      |
| `api`      | `string` | **Yes**  | Your secret API key for the chosen provider. **This is a secret and should not be committed.**                                                               |
| `cache`    | `object` | No       | Settings for caching AI responses to reduce costs and latency.                                                                                               |

### Provider-Specific Fields

These fields are only required if you are using the specified provider.

| Field                   | Provider            | Required | Description                                                             |
| :---------------------- | :------------------ | :------- | :---------------------------------------------------------------------- |
| `endpoint`              | `OpenAI Compatible` | **Yes**  | The base URL of the API endpoint (e.g. `https://api.example.com/v1`).   |
| `gcp_project_id`        | `Google Vertex AI`  | **Yes**  | Your Google Cloud Platform Project ID.                                  |
| `gcp_region`            | `Google Vertex AI`  | **Yes**  | The GCP region for your Vertex AI model (e.g. `us-central1`).           |
| `aws_region`            | `Amazon Bedrock`    | **Yes**  | The AWS region where your Bedrock models are hosted (e.g. `us-east-1`). |
| `aws_access_key_id`     | `Amazon Bedrock`    | **Yes**  | Your AWS Access Key ID for authentication.                              |
| `aws_secret_access_key` | `Amazon Bedrock`    | **Yes**  | Your AWS Secret Access Key for authentication.                          |
| `azure_resource_name`   | `Azure OpenAI`      | **Yes**  | The name of your Azure OpenAI resource.                                 |

### Custom Guidelines

| Field               | Type    | Required | Description                                                                                         |
| :------------------ | :------ | :------- | :-------------------------------------------------------------------------------------------------- |
| `commits.guides`    | `array` | No       | A list of paths to local `.md` or `.txt` files that will be used as guidelines for `gct ai commit`. |
| `changelogs.guides` | `array` | No       | A list of paths to local `.md` or `.txt` files that will be used as guidelines for `gct ai log`.    |

---

## Environment Variables

All configuration fields can be set using environment variables. This is especially useful in CI/CD environments.

| Environment Variable        | `gct.yaml` Field        | Required                              |
| :-------------------------- | :---------------------- | :------------------------------------ |
| `GCT_NAME`                  | `name`                  | No                                    |
| `GCT_PROVIDER`              | `provider`              | **Yes**                               |
| `GCT_MODEL`                 | `model`                 | **Yes**                               |
| `GCT_API_KEY`               | `api`                   | **Yes**                               |
| `GCT_ENDPOINT`              | `endpoint`              | Only for `OpenAI Compatible` provider |
| `GCT_GCP_PROJECT_ID`        | `gcp_project_id`        | Only for `Google Vertex AI` provider  |
| `GCT_GCP_REGION`            | `gcp_region`            | Only for `Google Vertex AI` provider  |
| `GCT_AWS_REGION`            | `aws_region`            | Only for `Amazon Bedrock` provider    |
| `GCT_AWS_ACCESS_KEY_ID`     | `aws_access_key_id`     | Only for `Amazon Bedrock` provider    |
| `GCT_AWS_SECRET_ACCESS_KEY` | `aws_secret_access_key` | Only for `Amazon Bedrock` provider    |
| `GCT_AZURE_RESOURCE_NAME`   | `azure_resource_name`   | Only for `Azure OpenAI` provider      |
| `GCT_CACHE_ENABLED`         | `cache.enabled`         | No                                    |
