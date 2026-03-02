# -- variables.tf
# ============================================================================
# Input Variables
# ============================================================================
# Configuration variables for Snowflake tasks
# ============================================================================

variable "task_configs" {
  description = "Map of configuration objects for Snowflake tasks"
  type = map(object({
    database      = string
    schema        = string
    name          = string
    warehouse     = optional(string, null)
    sql_statement = string

    # Schedule (for standalone/root tasks)
    schedule_minutes = optional(number, null)
    schedule_cron    = optional(string, null)

    # Optional settings
    comment                                  = optional(string, null)
    started                                  = optional(bool, false)
    allow_overlapping_execution              = optional(string, null)
    error_integration                        = optional(string, null)
    suspend_task_after_num_failures          = optional(number, null)
    user_task_timeout_ms                     = optional(number, null)
    user_task_managed_initial_warehouse_size = optional(string, null)

    # Task dependency (for DAG tasks)
    afters = optional(list(string), [])

    # Conditional execution
    when = optional(string, null)

    # Grants
    grants = optional(list(object({
      role_name  = string
      privileges = list(string)
    })), [])
  }))
  default = {}

  validation {
    condition     = alltrue([for k, task in var.task_configs : length(task.name) > 0])
    error_message = "Task name must not be empty."
  }

  validation {
    condition     = alltrue([for k, task in var.task_configs : length(task.database) > 0])
    error_message = "Database name must not be empty."
  }

  validation {
    condition     = alltrue([for k, task in var.task_configs : length(task.schema) > 0])
    error_message = "Schema name must not be empty."
  }

  validation {
    condition     = alltrue([for k, task in var.task_configs : length(task.sql_statement) > 0])
    error_message = "SQL statement must not be empty."
  }

  validation {
    condition = alltrue([
      for k, task in var.task_configs :
      task.user_task_managed_initial_warehouse_size == null ? true :
      contains(["XSMALL", "X-SMALL", "SMALL", "MEDIUM", "LARGE", "XLARGE", "X-LARGE", "XXLARGE", "X2LARGE", "2X-LARGE", "XXXLARGE", "X3LARGE", "3X-LARGE", "X4LARGE", "4X-LARGE", "X5LARGE", "5X-LARGE", "X6LARGE", "6X-LARGE"], upper(task.user_task_managed_initial_warehouse_size))
    ])
    error_message = "Invalid user_task_managed_initial_warehouse_size. Valid values: XSMALL, SMALL, MEDIUM, LARGE, XLARGE, XXLARGE, XXXLARGE, X4LARGE, X5LARGE, X6LARGE."
  }
}
