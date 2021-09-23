locals {
  app_name    = "github-quick-action"
  app_archive = "${var.app_binary_path}.zip"

  gateway_logformat_default = "{\"http_method\":\"$context.httpMethod\",\"path\":\"$context.path\",\"request_id\":\"$context.requestId\",\"lambda\":{\"status\":$context.integration.status,\"error\":\"$context.integration.error\"},\"response\":{\"status\":$context.status}}"
  gateway_logformat_verbose = "{\"time\":\"$context.requestTime\",\"protocol\":\"$context.protocol\",\"http_method\":\"$context.httpMethod\",\"gateway_api\":{\"id\":\"$context.apiId\",\"domain\":\"$context.domainName\",\"stage\":\"$context.stage\"},\"path\":\"$context.path\",\"request_id\":\"$context.requestId\",\"source_ip\":\"$context.identity.sourceIp\",\"user-agent\":\"$context.identity.userAgent\",\"lambda\":{\"lambda_status\":$context.integration.integrationStatus,\"status\":$context.integration.status,\"error\":\"$context.integration.error\",\"latency\":$context.integration.latency},\"response\":{\"status\":$context.status,\"latency\":$context.responseLatency,\"length\":$context.responseLength}}"
  gateway_logformat         = var.enable_tracing ? local.gateway_logformat_verbose : local.gateway_logformat_default
}

// Publish API Gateway as application ingress
module "api_gateway" {
  source = "terraform-aws-modules/apigateway-v2/aws"

  name          = "github-quick-actions-api"
  description   = "Github quick action application ingress"
  protocol_type = "HTTP"

  default_stage_access_log_destination_arn = module.app_lambda.lambda_cloudwatch_log_group_arn
  default_stage_access_log_format          = local.gateway_logformat
  default_route_settings = {
    detailed_metrics_enabled = true
    throttling_burst_limit   = 100
    throttling_rate_limit    = 100
  }

  integrations = {
    "POST /gh-webhook" = {
      lambda_arn             = module.app_lambda.lambda_function_arn
      payload_format_version = "2.0"
      timeout_milliseconds   = 1000 # NOTE: 1s timeout to avoid spamming
    }
  }

  tags = {
    "application.x-amz.com" : local.app_name
    "version.x-amz.com" : var.app_version
  }

  // NOTE: disable all unused features
  create_api_domain_name = false
  create_vpc_link        = false
}

// Publish lambda module using `terraform-aws-modules/lambda/aws`
module "app_lambda" {
  source = "terraform-aws-modules/lambda/aws"

  function_name = "github-quick-actions"
  description   = "Github quick action application"
  handler       = basename(var.app_binary_path)
  runtime       = "go1.x"

  create_package         = false
  publish                = true
  local_existing_package = data.archive_file.app_archive.output_path
  hash_extra             = data.archive_file.app_archive.output_sha

  environment_variables = {
    "GQA_GITHUB_APP_ID"         = var.github_app_id
    "GQA_GITHUB_PKEY"           = var.github_b64pkey
    "GQA_GITHUB_WEBHOOK_SECRET" = var.github_webhook_secret
    "GQA_LOG_LEVEL" : var.app_log_level
  }

  allowed_triggers = {
    AllowExecutionFromAPIGateway = {
      service    = "apigateway"
      source_arn = "${module.api_gateway.apigatewayv2_api_execution_arn}/*/*"
    },
  }

  attach_cloudwatch_logs_policy     = true
  cloudwatch_logs_retention_in_days = 1
  cloudwatch_logs_tags = {
    "application.x-amz.com" : local.app_name
    "version.x-amz.com" : var.app_version
  }

  attach_tracing_policy = var.enable_tracing
  tracing_mode          = var.enable_tracing ? "Active" : "PassThrough"

  tags = {
    "application.x-amz.com" : local.app_name
    "version.x-amz.com" : var.app_version
  }

  // NOTE: disable all unused features
  create_layer = false
}

data "archive_file" "app_archive" {
  type        = "zip"
  source_file = var.app_binary_path
  output_path = local.app_archive
}
