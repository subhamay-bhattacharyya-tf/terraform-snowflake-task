# -- examples/single-task/main.tf
# ============================================================================
# Single Snowflake Task Example
# ============================================================================
# This example demonstrates how to use the snowflake-task module
# to create a single Snowflake task.
# ============================================================================

module "task" {
  source = "../.."

  task_configs = var.task_configs
}
