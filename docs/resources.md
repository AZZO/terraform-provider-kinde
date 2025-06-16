# Kinde Provider Resources

This document provides detailed information about all resources available in the Kinde provider.

## Table of Contents

- [API](#api)
- [API Scope](#api-scope)
- [Application](#application)
- [Environment](#environment)
- [Permission](#permission)
- [Role](#role)

## API

Manages a Kinde API.

### Example Usage

```hcl
resource "kinde_api" "example" {
  name     = "example-api"
  audience = "https://api.example.com"
}
```

### Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `name` | `string` | Yes | The name of the API. Must be between 1 and 64 characters. |
| `audience` | `string` | Yes | The audience for the API. Must be between 1 and 64 characters. |

### Attributes

| Name | Type | Description |
|------|------|-------------|
| `api_id` | `string` | Unique identifier for the API. |

### Notes

- APIs are immutable and cannot be updated. Any changes will result in the API being replaced.
- The `api_id` is automatically generated and cannot be set manually.

## API Scope

Manages a scope for a Kinde API.

### Example Usage

```hcl
resource "kinde_api_scope" "example" {
  api_id       = kinde_api.example.api_id
  name         = "read:data"
  description  = "Read access to data"
}
```

### Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `api_id` | `string` | Yes | The ID of the API this scope belongs to. |
| `name` | `string` | Yes | The name of the scope. |
| `description` | `string` | No | A description of what the scope allows. |

### Attributes

| Name | Type | Description |
|------|------|-------------|
| `scope_id` | `string` | Unique identifier for the scope. |

### Notes

- The `scope_id` is automatically generated and cannot be set manually.
- Scopes are unique within an API.

## Application

Manages a Kinde application.

### Example Usage

```hcl
resource "kinde_application" "example" {
  name = "example-app"
  type = "regular_web"
}
```

### Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `name` | `string` | Yes | The name of the application. |
| `type` | `string` | Yes | The type of application. Valid values are: `regular_web`, `single_page`, `native`, `machine_to_machine`. |

### Attributes

| Name | Type | Description |
|------|------|-------------|
| `id` | `string` | Unique identifier for the application. |
| `client_id` | `string` | The client ID for the application. |
| `client_secret` | `string` | The client secret for the application. |

### Notes

- The `client_id` and `client_secret` are automatically generated and cannot be set manually.
- The `client_secret` is only available when the application is created and cannot be retrieved later.

## Environment

Manages a Kinde environment.

### Example Usage

```hcl
resource "kinde_environment" "example" {
  name = "example-environment"
}
```

### Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `name` | `string` | Yes | The name of the environment. |

### Attributes

| Name | Type | Description |
|------|------|-------------|
| `code` | `string` | Unique identifier for the environment. |
| `is_default` | `bool` | Whether this is the default environment. |
| `is_live` | `bool` | Whether this is a live environment. |
| `kinde_domain` | `string` | The Kinde domain for the environment. |
| `custom_domain` | `string` | The custom domain for the environment, if set. |
| `logo` | `string` | URL to the environment's logo. |
| `logo_dark` | `string` | URL to the environment's dark mode logo. |
| `favicon_svg` | `string` | URL to the environment's SVG favicon. |
| `favicon_fallback` | `string` | URL to the environment's fallback favicon. |
| `created_on` | `string` | When the environment was created. |

### Notes

- The `code` is automatically generated and cannot be set manually.
- Most attributes are read-only and cannot be modified.

## Permission

Manages a Kinde permission.

### Example Usage

```hcl
resource "kinde_permission" "example" {
  name        = "read:data"
  key         = "read:data"
  description = "Permission to read data"
}
```

### Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `name` | `string` | Yes | The name of the permission. |
| `key` | `string` | Yes | The key for the permission. |
| `description` | `string` | No | A description of what the permission allows. |

### Attributes

| Name | Type | Description |
|------|------|-------------|
| `id` | `string` | Unique identifier for the permission. |

### Notes

- The `id` is automatically generated and cannot be set manually.
- Permissions are unique within an organization.

## Role

The `kinde_role` resource manages a Kinde role.

### Example Usage

```hcl
resource "kinde_role" "example" {
  name          = "Admin"
  description   = "Administrator role"
  key           = "admin"
  is_default_role = false
  permissions   = ["read:users", "write:users"]
  scopes        = ["read:data", "write:data"]
}
```

### Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `role_id` | `string` | No | Unique identifier for the role. Automatically generated if not provided. |
| `name` | `string` | Yes | The name of the role. |
| `description` | `string` | No | The description of the role. |
| `key` | `string` | Yes | The key of the role. Must be unique. |
| `is_default_role` | `bool` | No | Whether this is a default role. Defaults to false. |
| `permissions` | `list(string)` | No | List of permission IDs associated with this role. |
| `scopes` | `list(string)` | No | List of scope IDs associated with this role. |

### Attributes

| Name | Type | Description |
|------|------|-------------|
| `role_id` | `string` | Unique identifier for the role. |
| `name` | `string` | The name of the role. |
| `description` | `string` | The description of the role. |
| `key` | `string` | The key of the role. |
| `is_default_role` | `bool` | Whether this is a default role. |
| `permissions` | `list(string)` | List of permission IDs associated with this role. |
| `scopes` | `list(string)` | List of scope IDs associated with this role. |

### Notes

- The `role_id` is automatically generated if not provided.
- The `key` must be unique and cannot be changed after creation.
- The `permissions` and `scopes` lists are optional and can be updated at any time.
- When updating `permissions` or `scopes`, the provider will automatically handle adding and removing items to match the desired state. 