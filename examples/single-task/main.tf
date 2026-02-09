# Example: Single Snowflake Task
#
# This example demonstrates how to use the snowflake-task module
# to create a single Snowflake task.

module "task" {
  source = "../../modules/snowflake-task"

  task_configs = var.task_configs
}
