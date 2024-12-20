# Edge Leap Configuration Schema v1.0

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
| `target-condition` | string | Condition for deployment targeting (when in `draft` mode this is set automatically) |

### `device`
Device identification details.

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Unique device identifier in the IoT Hub|

### `infra`
Infrastructure configuration.

| Field | Type | Description |
|-------|------|-------------|
| `hub` | string | Name of the IoT Hub where the target device is connected |

### `module`
Module-specific configuration.

| Field | Type | Description |
|-------|------|-------------|
| `create-options` | string | Docker container creation options in JSON format |
| `env` | array of strings | Environment variables for the module in the format `"MY_VAR=MY_VAL"` |
| `image` | string | Docker image reference |
| `name` | string | Name of the module |
| `startup-order` | integer | Startup sequence priority |

## Automatically Generated Fields
This section is automatically generated by the tool and should not be modified.

| Field | Type | Description |
|-------|------|-------------|
| `session` | string | Session identifier |
| `version` | integer | Configuration schema version |