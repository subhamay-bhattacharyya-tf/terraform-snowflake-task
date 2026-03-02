# Terraform Snowflake Module - Task

![Release](https://github.com/subhamay-bhattacharyya-tf/terraform-snowflake-task/actions/workflows/ci.yaml/badge.svg)&nbsp;![Snowflake](https://img.shields.io/badge/Snowflake-29B5E8?logo=snowflake&logoColor=white)&nbsp;![Commit Activity](https://img.shields.io/github/commit-activity/t/subhamay-bhattacharyya-tf/terraform-snowflake-task)&nbsp;![Last Commit](https://img.shields.io/github/last-commit/subhamay-bhattacharyya-tf/terraform-snowflake-task)&nbsp;![Release Date](https://img.shields.io/github/release-date/subhamay-bhattacharyya-tf/terraform-snowflake-task)&nbsp;![Repo Size](https://img.shields.io/github/repo-size/subhamay-bhattacharyya-tf/terraform-snowflake-task)&nbsp;![File Count](https://img.shields.io/github/directory-file-count/subhamay-bhattacharyya-tf/terraform-snowflake-task)&nbsp;![Issues](https://img.shields.io/github/issues/subhamay-bhattacharyya-tf/terraform-snowflake-task)&nbsp;![Top Language](https://img.shields.io/github/languages/top/subhamay-bhattacharyya-tf/terraform-snowflake-task)&nbsp;![Custom Endpoint](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/bsubhamay/3250f218c58c0b8262fbeff9dea5f412/raw/terraform-snowflake-task.json?)

A Terraform module for creating and managing Snowflake tasks using a map of configuration objects. Supports creating single tasks or task DAGs (Directed Acyclic Graphs) with a single module call.

## Features

- Map-based configuration for creating single or multiple tasks
- Support for task DAGs with `afters` dependencies
- Role-based access control with `grants` configuration
- Built-in input validation with descriptive error messages
- Sensible defaults for optional properties
- Outputs keyed by task identifier for easy reference
- Support for scheduled tasks (cron or interval-based)
- Support for serverless tasks (user-managed warehouse size)
- Conditional task execution with `when` clause

## Usage

### Single Task

```hcl
module "task" {
  source = "github.com/subhamay-bhattacharyya-tf/terraform-snowflake-task"

  task_configs = {
    "my_task" = {
      database         = "MY_DATABASE"
      schema           = "MY_SCHEMA"
      name             = "MY_TASK"
      warehouse        = "MY_WAREHOUSE"
      sql_statement    = "CALL my_procedure()"
      schedule_minutes = 60
      started          = false
      comment          = "My scheduled task"
      grants = [
        {
          role_name  = "DATA_ENGINEER"
          privileges = ["MONITOR", "OPERATE"]
        }
      ]
    }
  }
}
```

### Task DAG (Multiple Tasks)

```hcl
locals {
  tasks = {
    "root_task" = {
      database         = "MY_DATABASE"
      schema           = "MY_SCHEMA"
      name             = "ROOT_TASK"
      warehouse        = "MY_WAREHOUSE"
      sql_statement    = "CALL stage_data()"
      schedule_minutes = 60
      started          = false
      comment          = "Root task - runs every hour"
      grants = [
        {
          role_name  = "DATA_ENGINEER"
          privileges = ["MONITOR", "OPERATE"]
        }
      ]
    }
    "transform_task" = {
      database      = "MY_DATABASE"
      schema        = "MY_SCHEMA"
      name          = "TRANSFORM_TASK"
      warehouse     = "MY_WAREHOUSE"
      sql_statement = "CALL transform_data()"
      afters        = ["MY_DATABASE.MY_SCHEMA.ROOT_TASK"]
      started       = false
      comment       = "Transform task - runs after root"
      grants = [
        {
          role_name  = "DATA_ENGINEER"
          privileges = ["MONITOR"]
        }
      ]
    }
    "load_task" = {
      database      = "MY_DATABASE"
      schema        = "MY_SCHEMA"
      name          = "LOAD_TASK"
      warehouse     = "MY_WAREHOUSE"
      sql_statement = "CALL load_to_final()"
      afters        = ["MY_DATABASE.MY_SCHEMA.TRANSFORM_TASK"]
      started       = false
      comment       = "Load task - runs after transform"
      grants = [
        {
          role_name  = "DATA_ENGINEER"
          privileges = ["MONITOR"]
        }
      ]
    }
  }
}

module "tasks" {
  source = "github.com/subhamay-bhattacharyya-tf/terraform-snowflake-task"

  task_configs = local.tasks
}
```

### Serverless Task

```hcl
module "serverless_task" {
  source = "github.com/subhamay-bhattacharyya-tf/terraform-snowflake-task"

  task_configs = {
    "serverless_task" = {
      database         = "MY_DATABASE"
      schema           = "MY_SCHEMA"
      name             = "SERVERLESS_TASK"
      sql_statement    = "SELECT 1"
      schedule_minutes = 5
      user_task_managed_initial_warehouse_size = "XSMALL"
      started          = false
      comment          = "Serverless task using managed compute"
    }
  }
}
```

## Examples

- [Single Task](examples/single-task) - Create a single scheduled task
- [Multiple Tasks (DAG)](examples/multiple-tasks) - Create a task DAG with dependencies

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.3.0 |
| snowflake | >= 1.0.0 |

## Providers

| Name | Version |
|------|---------|
| snowflake (snowflakedb/snowflake) | >= 1.0.0 |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|----------|
| task_configs | Map of configuration objects for Snowflake tasks | `map(object)` | `{}` | no |

### task_configs Object Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| database | string | - | Database where the task resides (required) |
| schema | string | - | Schema where the task resides (required) |
| name | string | - | Task identifier (required) |
| warehouse | string | null | Warehouse to use for task execution |
| sql_statement | string | - | SQL statement to execute (required) |
| schedule_minutes | number | null | Schedule interval in minutes (for root tasks) |
| schedule_cron | string | null | Cron expression for scheduling (for root tasks) |
| comment | string | null | Description of the task |
| started | bool | false | Whether the task is started |
| allow_overlapping_execution | string | null | Allow concurrent task runs ("true"/"false") |
| error_integration | string | null | Error notification integration name |
| suspend_task_after_num_failures | number | null | Suspend after N consecutive failures |
| user_task_timeout_ms | number | null | Task timeout in milliseconds |
| user_task_managed_initial_warehouse_size | string | null | Warehouse size for serverless tasks |
| afters | list(string) | [] | List of predecessor task fully qualified names (database.schema.name) |
| when | string | null | Conditional execution expression |
| grants | list(object) | [] | List of role grants with privileges |

### grants Object Properties

| Property | Type | Description |
|----------|------|-------------|
| role_name | string | Name of the role to grant privileges to |
| privileges | list(string) | List of privileges to grant (e.g., MONITOR, OPERATE) |

### Valid Warehouse Sizes (for serverless tasks)

- XSMALL (X-SMALL)
- SMALL
- MEDIUM
- LARGE
- XLARGE (X-LARGE)
- XXLARGE (2X-LARGE)
- XXXLARGE (3X-LARGE)
- X4LARGE (4X-LARGE)
- X5LARGE (5X-LARGE)
- X6LARGE (6X-LARGE)

## Outputs

| Name | Description |
|------|-------------|
| task_names | Map of task names keyed by identifier |
| task_fully_qualified_names | Map of fully qualified task names |
| task_databases | Map of task databases |
| task_schemas | Map of task schemas |
| task_states | Map of task states (started or suspended) |
| tasks | All task resources |

## Validation

The module validates inputs and provides descriptive error messages for:

- Empty task name
- Empty database name
- Empty schema name
- Empty SQL statement
- Invalid warehouse size for serverless tasks
- Negative suspend_task_after_num_failures value

## Testing

The module includes Terratest-based integration tests:

```bash
cd test
go mod tidy
go test -v -timeout 30m
```

Required environment variables for testing:
- `SNOWFLAKE_ORGANIZATION_NAME` - Snowflake organization name
- `SNOWFLAKE_ACCOUNT_NAME` - Snowflake account name
- `SNOWFLAKE_USER` - Snowflake username
- `SNOWFLAKE_ROLE` - Snowflake role (e.g., "SYSADMIN")
- `SNOWFLAKE_PRIVATE_KEY` - Snowflake private key for key-pair authentication

## CI/CD Configuration

The CI workflow runs on:
- Push to `main`, `feature/**`, and `bug/**` branches (when `*.tf` or `examples/**` changes)
- Pull requests to `main` (when `*.tf` or `examples/**` changes)
- Manual workflow dispatch

The workflow includes:
- Terraform validation and format checking
- Examples validation
- Terratest integration tests (output displayed in GitHub Step Summary)
- Changelog generation (non-main branches)
- Semantic release (main branch only)

The CI workflow uses the following GitHub organization variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `TERRAFORM_VERSION` | Terraform version for CI jobs | `1.3.0` |
| `GO_VERSION` | Go version for Terratest | `1.21` |
| `SNOWFLAKE_ORGANIZATION_NAME` | Snowflake organization name | - |
| `SNOWFLAKE_ACCOUNT_NAME` | Snowflake account name | - |
| `SNOWFLAKE_USER` | Snowflake username | - |
| `SNOWFLAKE_ROLE` | Snowflake role (e.g., SYSADMIN) | - |

The following GitHub secrets are required for Terratest integration tests:

| Secret | Description | Required |
|--------|-------------|----------|
| `SNOWFLAKE_PRIVATE_KEY` | Snowflake private key for key-pair authentication | Yes |

## License

MIT License - See [LICENSE](LICENSE) for details.
