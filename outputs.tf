# -- outputs.tf
# ============================================================================
# Output Values
# ============================================================================
# Outputs for Snowflake tasks
# ============================================================================

output "task_names" {
  description = "The names of the created tasks."
  value       = { for k, v in snowflake_task.this : k => v.name }
}

output "task_fully_qualified_names" {
  description = "The fully qualified names of the tasks."
  value       = { for k, v in snowflake_task.this : k => v.fully_qualified_name }
}

output "task_databases" {
  description = "The databases of the tasks."
  value       = { for k, v in snowflake_task.this : k => v.database }
}

output "task_schemas" {
  description = "The schemas of the tasks."
  value       = { for k, v in snowflake_task.this : k => v.schema }
}

output "task_states" {
  description = "The states of the tasks (started or suspended)."
  value       = { for k, task in var.task_configs : k => task.started ? "started" : "suspended" }
}

output "tasks" {
  description = "All task resources."
  value       = snowflake_task.this
}
