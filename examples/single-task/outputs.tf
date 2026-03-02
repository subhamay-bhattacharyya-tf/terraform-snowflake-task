# -- examples/single-task/outputs.tf
# ============================================================================
# Output Values
# ============================================================================
# Outputs for the single task example
# ============================================================================

output "task_names" {
  description = "The names of the created tasks"
  value       = module.task.task_names
}

output "task_fully_qualified_names" {
  description = "The fully qualified names of the tasks"
  value       = module.task.task_fully_qualified_names
}

output "task_databases" {
  description = "The databases of the tasks"
  value       = module.task.task_databases
}

output "task_schemas" {
  description = "The schemas of the tasks"
  value       = module.task.task_schemas
}

output "task_states" {
  description = "The states of the tasks (started or suspended)"
  value       = module.task.task_states
}
