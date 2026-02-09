# Basic Example - Single Task

This example demonstrates how to create a single Snowflake task using the module.

## Usage

```hcl
module "task" {
  source = "../../modules/snowflake-task"

  task_configs = {
    "my_task" = {
      database         = "MY_DATABASE"
      schema           = "MY_SCHEMA"
      name             = "MY_TASK"
      warehouse        = "MY_WAREHOUSE"
      sql_statement    = "SELECT 1"
      schedule_minutes = 60
      enabled          = false
      comment          = "My test task"
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
