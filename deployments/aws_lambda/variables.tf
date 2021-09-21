// Package information
variable "app_binary_path" {
  description = "Application compiled binary path."
  type        = string

  validation {
    condition     = fileexists(var.app_binary_path)
    error_message = "The application binary must exists."
  }
}
variable "app_version" {
  description = "Application version."
  type        = string

  validation {
    # Source: https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
    condition     = can(regex("^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$", var.app_version))
    error_message = "The application version only accepts SemVer compatible strings."
  }
}

variable "enable_tracing" {
  description = "Enable 'tracing' mode; it enables AWS X-Ray and add more verbosity to Gateway logs."
  type        = bool
  default     = false
}

// Github Application information
variable "github_app_id" {
  description = "Github application ID."
  sensitive   = true
  type        = string
}
variable "github_b64pkey" {
  description = "Github application base64 encoded private key."
  sensitive   = true
  type        = string
}
variable "github_webhook_secret" {
  description = "Github application webhook secret."
  sensitive   = true
  type        = string
}

variable "app_log_level" {
  description = "Application log level."
  type        = string
  default     = "info"
}
