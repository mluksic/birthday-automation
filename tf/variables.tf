variable "APP_ENV" {
  description = "app env"
  type        = string
  sensitive   = false
}

variable "ACCOUNT_SID" {
  description = "Account is for SMS messaging"
  type        = string
  sensitive   = true
}

variable "AUTH_TOKEN" {
  description = "Auth token for SMS messaging service"
  type        = string
  sensitive   = true
}

variable "PHONE_NUMBER" {
  description = "Phone number for SMS"
  type        = string
  sensitive   = true
}

variable "SERVICE_ID" {
  description = "Service ID of SMS service"
  type        = string
  sensitive   = true
}

variable "APP_PASSWORD" {
  description = "APP password"
  type        = string
  sensitive   = true
}
variable "FROM_EMAIL" {
  description = "Email from where the mail will be sent to"
  type        = string
  sensitive   = true
}

variable "EMAIL_RECEIVERS" {
  description = "List of emails that will receive the email"
  type        = string
  sensitive   = true
}
