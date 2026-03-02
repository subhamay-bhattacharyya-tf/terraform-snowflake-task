# Multiple Tasks Example - Task DAG

This example demonstrates how to create multiple Snowflake tasks forming a DAG (Directed Acyclic Graph) using the module.

## Usage

```hcl
module "tasks" {
  source = "../.."

  task_configs = {
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
      afters        = ["ROOT_TASK"]
      started       = false
      comment       = "Transform task - runs after root task"
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
      afters        = ["TRANSFORM_TASK"]
      started       = false
      comment       = "Load task - runs after transform task"
      grants = [
        {
          role_name  = "DATA_ENGINEER"
          privileges = ["MONITOR"]
        }
      ]
    }
  }
}
```

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| task_configs | Map of task configuration objects | map(object) | yes |

## Outputs

| Name | Description |
|------|-------------|
| task_names | The names of the created tasks |
| task_fully_qualified_names | The fully qualified names of the tasks |
| task_databases | The databases of the tasks |
| task_schemas | The schemas of the tasks |
| task_states | The states of the tasks |
| tasks | All task resources |

## Task DAG Structure

```
ROOT_TASK (scheduled: every 60 minutes)
    │
    ▼
TRANSFORM_TASK (triggered by ROOT_TASK)
    │
    ▼
LOAD_TASK (triggered by TRANSFORM_TASK)
```
