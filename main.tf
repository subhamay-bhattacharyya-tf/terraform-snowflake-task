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

  # Task dependency (for DAG tasks) - conflicts with schedule
  after = length(each.value.afters) > 0 ? each.value.afters : null

  # Conditional execution
  when = each.value.when
}

# Grant privileges on tasks to roles
resource "snowflake_grant_privileges_to_account_role" "task_grants" {
  for_each = {
    for grant in flatten([
      for task_key, task in var.task_configs : [
        for grant in task.grants : {
          key        = "${task_key}_${grant.role_name}"
          task_key   = task_key
          database   = task.database
          schema     = task.schema
          task_name  = task.name
          role_name  = grant.role_name
          privileges = grant.privileges
        }
      ]
    ]) : grant.key => grant
  }

  privileges        = each.value.privileges
  account_role_name = each.value.role_name
  on_schema_object {
    object_type = "TASK"
    object_name = "\"${each.value.database}\".\"${each.value.schema}\".\"${each.value.task_name}\""
  }

  depends_on = [snowflake_task.this]
}
