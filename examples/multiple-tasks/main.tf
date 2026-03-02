# Example: Multiple Snowflake Tasks (DAG)
#
# This example demonstrates how to use the snowflake-task module
# to create multiple Snowflake tasks forming a DAG (Directed Acyclic Graph).

module "tasks" {
  source = "../.."

  task_configs = var.task_configs
}
