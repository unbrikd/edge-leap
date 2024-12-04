# Application Configuration Schema

## Overview
This configuration defines the settings for a deployment configuration, likely for an IoT or microservices application.

## Configuration Sections

### `auth`
Defines authentication credentials.

| Field | Type | Description | 
|-------|------|-------------|
| `token` | string | Shared Access Signature (SAS) token for authentication |

### `deployment`
Deployment-specific configuration.

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier for the deployment |
| `priority` | integer | Deployment priority level |
| `target-condition` | string | Condition for deployment targeting (currently empty) |

### `device`
Device identification details.

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Unique device identifier |

### `infra`
Infrastructure configuration.

| Field | Type | Description |
|-------|------|-------------|
| `hub` | string | Name of the IoT hub |

### `module`
Module-specific configuration.

| Field | Type | Description |
|-------|------|-------------|
| `create-options` | string | Docker container creation options in JSON format |
| `env` | array of strings | Environment variables for the module |
| `image` | string | Docker image reference |
| `name` | string | Name of the module |
| `startup-order` | integer | Startup sequence priority |

### Additional Configuration

| Field | Type | Description |
|-------|------|-------------|
| `session` | string | Session identifier |
| `version` | integer | Configuration version |