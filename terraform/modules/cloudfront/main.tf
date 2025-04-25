resource "aws_cloudfront_distribution" "this" {
  enabled             = true
  is_ipv6_enabled     = true
  price_class         = var.price_class
  aliases             = ["${var.name}.${var.domain_name}"]
  wait_for_deployment = var.wait_for_deployment

  origin {
    domain_name = var.origin_domain_name
    origin_id   = "${var.name}-origin"

    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = var.origin_protocol_policy
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  default_cache_behavior {
    allowed_methods  = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "${var.name}-origin"

    forwarded_values {
      query_string = true
      cookies {
        forward = "all"
      }
      headers = var.forwarded_headers
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = var.cache_policy.min_ttl
    default_ttl            = var.cache_policy.default_ttl
    max_ttl                = var.cache_policy.max_ttl
  }

  viewer_certificate {
    acm_certificate_arn      = var.certificate_arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  tags = merge(
    var.shared_tags,
    {
      "TF"      = "true",
      "Project" = var.project,
      "Github"  = "github.com/Mad-Pixels/applingo-api",
      "Name"    = "${var.name}-cdn"
    }
  )
}