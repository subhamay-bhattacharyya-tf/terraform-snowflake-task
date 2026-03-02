# Example: Single Snowflake Task
#
# This example demonstrates how to use the snowflake-task module
# to create a single Snowflake task.

module "task" {
  source = "../.."

  task_configs = var.task_configs
}
