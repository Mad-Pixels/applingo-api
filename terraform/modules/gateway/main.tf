resource "aws_api_gateway_rest_api" "this" {
  name = "${var.project}-${var.api_name}"
  body = local.manifest

  endpoint_configuration {
    types = ["REGIONAL"]
  }

  tags = merge(
    var.shared_tags,
    {
      "Name" = "${var.project}-${var.api_name}",
    }
  )
}

resource "aws_api_gateway_deployment" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id

  triggers = {
    redeployment = sha1(jsonencode(aws_api_gateway_rest_api.this.body))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "this" {
  deployment_id = aws_api_gateway_deployment.this.id
  rest_api_id   = aws_api_gateway_rest_api.this.id
  stage_name    = var.stage_name
  tags          = var.shared_tags

  xray_tracing_enabled = false

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.this.arn
    format          = "{\"requestId\":\"$context.requestId\",\"traceId\":\"$context.xrayTraceId\",\"ip\":\"$context.identity.sourceIp\",\"caller\":\"$context.identity.caller\",\"user\":\"$context.identity.user\",\"requestTime\":\"$context.requestTime\",\"httpMethod\":\"$context.httpMethod\",\"domainName\":\"$context.domainName\",\"resourcePath\":\"$context.resourcePath\",\"status\":\"$context.status\",\"error\":\"$context.error.message\",\"validationError\":\"$context.error.validationErrorString\",\"lambdaStatus\":\"$context.integration.status\",\"lambdaLatency\":\"$context.integration.latency\",\"authStatus\":\"$context.authorizer.status\",\"protocol\":\"$context.protocol\",\"responseLength\":\"$context.responseLength\",\"user-agent\":\"$context.identity.userAgent\",\"wafResponse\":\"$context.wafResponseCode\",\"wafError\":\"$context.waf.error\"}"
  }

  depends_on = [aws_api_gateway_account.this]
}

resource "aws_api_gateway_method_settings" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  stage_name  = aws_api_gateway_stage.this.stage_name
  method_path = "*/*"

  settings {
    metrics_enabled    = true
    data_trace_enabled = false
    logging_level      = "ERROR"

    throttling_rate_limit  = 1000
    throttling_burst_limit = 500
  }
}

resource "aws_api_gateway_domain_name" "this" {
  count = var.custom_domain != "" ? 1 : 0

  domain_name              = var.custom_domain
  regional_certificate_arn = var.acm_certificate_arn
  security_policy          = "TLS_1_2"

  endpoint_configuration {
    types = ["REGIONAL"]
  }

  tags = var.shared_tags
}

resource "aws_api_gateway_base_path_mapping" "this" {
  count = var.custom_domain != "" ? 1 : 0

  api_id      = aws_api_gateway_rest_api.this.id
  stage_name  = aws_api_gateway_stage.this.stage_name
  domain_name = aws_api_gateway_domain_name.this[0].domain_name
}

resource "aws_wafv2_web_acl_association" "this" {
  count = var.wafv2_web_acl_arn != "" ? 1 : 0

  resource_arn = aws_api_gateway_stage.this.arn
  web_acl_arn  = var.wafv2_web_acl_arn
}

resource "aws_api_gateway_account" "this" {
  cloudwatch_role_arn = aws_iam_role.cloudwatch_role.arn
}