resource "aws_alb" "main" {

    internal = false
    load_balancer_type = "application"
    security_groups = [ aws_security_group.lb.id ]
    subnets = aws_subnet.public[*].id
     
    tags = {
        name = "${var.env}-aoc2024-lb"
        environment = var.env
    }
}

resource "aws_alb_target_group" "app" {
    vpc_id = aws_vpc.vpc.id
    port = var.app_tcp_port
    protocol = "HTTP" 
    target_type = "ip"

    health_check {
        path = var.health_check_path
        protocol = "HTTP"
        port = "traffic-port"
        matcher = "200"
    }
    tags = {
        name = "${var.env}-aoc2024-lb-grp"
        environment = var.env
    }
}

resource "aws_alb_listener" "frontend_http" {
  count = var.alb_http ? 1 : 0
  load_balancer_arn = aws_alb.main.arn
  port = var.alb_http_port
  protocol = "HTTP"

  default_action {
    type = "forward"
    target_group_arn = aws_alb_target_group.app.arn
  }
}

resource "aws_alb_listener" "frontend_http_redirect" {
  count = var.alb_https && !var.alb_http ? 1 : 0
  load_balancer_arn = aws_alb.main.arn
  port = var.alb_http_port
  protocol = "HTTP"

  default_action {
    type = "redirect"
    redirect {
      protocol    = "HTTPS"
      port        = "443"
      status_code = "HTTP_301"
    }
  }
}
resource "aws_lb_listener" "frontend_https" {
    count = var.alb_https ? 1 : 0

    load_balancer_arn = aws_alb.main.arn
    port              = var.alb_https_port
    protocol          = "HTTPS"

    ssl_policy      = "ELBSecurityPolicy-TLS13-1-0-2021-06"

    // to make sure cert is validated before moving on
    depends_on = [aws_acm_certificate_validation.cert[0]]

    certificate_arn = aws_acm_certificate.alb_cert[0].arn

    default_action {
        type = "forward"
        target_group_arn = aws_alb_target_group.app.id
    }
}

// dns name registration
data "aws_route53_zone" "primary" {
    name = var.dns_zone
}

resource "aws_route53_record" "app_alias" {
    zone_id = data.aws_route53_zone.primary.zone_id
    name = var.alb_dns_name
    type = "A"

    alias {
        name = aws_alb.main.dns_name
        zone_id = aws_alb.main.zone_id
        evaluate_target_health = true
    }
}

// https for front end
resource "aws_acm_certificate" "alb_cert" {
    count = var.alb_https ? 1 : 0
    
    domain_name = var.alb_dns_name
    validation_method = "DNS"
    tags = {
        name = "${var.env}-aoc2024-alb-cert"
        environment = var.env
    }
}

resource "aws_route53_record" "cert_validation" {
    for_each = var.alb_https ? {
        for dvo in aws_acm_certificate.alb_cert[0].domain_validation_options : dvo.domain_name => {
          name  = dvo.resource_record_name
          type  = dvo.resource_record_type
          value = dvo.resource_record_value
        }
    } : {}

    zone_id = data.aws_route53_zone.primary.zone_id
    name    = each.value.name
    type    = each.value.type
    ttl     = 60
    records = [each.value.value]
}

resource "aws_acm_certificate_validation" "cert" {
    count = var.alb_https ? 1 : 0

    certificate_arn         = aws_acm_certificate.alb_cert[0].arn
    validation_record_fqdns = [for record in aws_route53_record.cert_validation : record.fqdn]
}