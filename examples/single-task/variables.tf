variable "task_configs" {
  description = "Map of configuration objects for Snowflake tasks"
  type = map(object({
    database      = string
    schema        = string
    name          = string
    warehouse     = optional(string, null)
    sql_statement = string

    schedule_minutes = optional(number, null)
    schedule_cron    = optional(string, null)

    comment                                  = optional(string, null)
    started                                  = optional(bool, false)
    allow_overlapping_execution              = optional(string, null)
    error_integration                        = optional(string, null)
    suspend_task_after_num_failures          = optional(number, null)
    user_task_timeout_ms                     = optional(number, null)
    user_task_managed_initial_warehouse_size = optional(string, null)

    afters = optional(list(string), [])
    when   = optional(string, null)
  }))
  default = {}
}

# Snowflake authentication variables
variable "snowflake_organization_name" {
  description = "Snowflake organization name"
  type        = string
  default     = null
}

variable "snowflake_account_name" {
  description = "Snowflake account name"
  type        = string
  default     = null
}

variable "snowflake_user" {
  description = "Snowflake username"
  type        = string
  default     = null
}

variable "snowflake_role" {
  description = "Snowflake role"
  type        = string
  default     = null
}

variable "snowflake_private_key" {
  description = "Snowflake private key for key-pair authentication"
  type        = string
  sensitive   = true
  default     = null
}
