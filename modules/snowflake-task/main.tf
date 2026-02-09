# Snowflake Task Resource
# Creates and manages one or more Snowflake tasks based on the task_configs map

resource "snowflake_task" "this" {
  for_each = var.task_configs

  database      = each.value.database
  schema        = each.value.schema
  name          = each.value.name
  warehouse     = each.value.warehouse
  sql_statement = each.value.sql_statement
  started       = each.value.started

  # Schedule configuration (standalone tasks only)
  dynamic "schedule" {
    for_each = each.value.schedule_minutes != null || each.value.schedule_cron != null ? [1] : []
    content {
      minutes    = each.value.schedule_minutes
      using_cron = each.value.schedule_cron
    }
  }

  # Optional settings
  comment                                  = each.value.comment
  allow_overlapping_execution              = each.value.allow_overlapping_execution
  error_integration                        = each.value.error_integration
  suspend_task_after_num_failures          = each.value.suspend_task_after_num_failures
  user_task_timeout_ms                     = each.value.user_task_timeout_ms
  user_task_managed_initial_warehouse_size = each.value.user_task_managed_initial_warehouse_size

  # Task dependency (for DAG tasks)
  after = each.value.afters

  # Conditional execution
  when = each.value.when
}
