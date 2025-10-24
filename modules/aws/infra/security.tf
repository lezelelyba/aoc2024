locals {
  lb_ingress_rules = [
    {
      protocol    = "tcp"
      from_port   = tonumber(var.alb_http_port)
      to_port     = tonumber(var.alb_http_port)
      cidr_blocks = ["0.0.0.0/0"]
    }
  ]
  # append HTTPS rule if HTTPS is enabled
  lb_ingress_rules_https = var.alb_https ? [
    {
      protocol    = "tcp"
      from_port   = tonumber(var.alb_https_port)
      to_port     = tonumber(var.alb_https_port)
      cidr_blocks = ["0.0.0.0/0"]
    }
  ] : []
  
  combined_rules = concat(local.lb_ingress_rules, local.lb_ingress_rules_https)
}
resource "aws_security_group" "lb" {
    vpc_id = aws_vpc.vpc.id
    name = "${var.env}-aoc2024-lb-sg"

    tags = {
        name = "${var.env}-aoc2024-lb-sg"
        environment = var.env
    }

    dynamic "ingress" {
        for_each = local.combined_rules
        content {
          protocol    = ingress.value.protocol
          from_port   = ingress.value.from_port
          to_port     = ingress.value.to_port
          cidr_blocks = ingress.value.cidr_blocks
        }
    }


    egress {
        protocol = "-1"
        from_port = 0
        to_port = 0
        cidr_blocks = ["0.0.0.0/0"]
    }
}

resource "aws_security_group" "ecs" {
    vpc_id = aws_vpc.vpc.id
    name = "${var.env}-aoc2024-ecs-sg"

    tags = {
        name = "${var.env}-aoc2024-ecs-sg"
        environment = var.env
    }

    ingress {
        protocol = "tcp"
        from_port = var.app_tcp_port
        to_port = var.app_tcp_port
        security_groups = [aws_security_group.lb.id]
    }

    egress {
        protocol = "-1"
        from_port = 0
        to_port = 0
        cidr_blocks = ["0.0.0.0/0"]
    }
}
resource "aws_security_group" "bastion" {
    vpc_id = aws_vpc.vpc.id
    name = "${var.env}-aoc2024-bastion-sg"

    tags = {
        name = "${var.env}-aoc2024-bastion-sg"
        environment = var.env
    }

    ingress {
        protocol = "tcp"
        from_port = 22
        to_port = 22
        cidr_blocks = ["0.0.0.0/0"]
    }

    egress {
        protocol = "-1"
        from_port = 0
        to_port = 0
        cidr_blocks = ["0.0.0.0/0"]
    }
}
resource "aws_security_group" "testhost" {
    vpc_id = aws_vpc.vpc.id
    name = "${var.env}-aoc2024-testhost-sg"

    tags = {
        name = "${var.env}-aoc2024-testhost-sg"
        environment = var.env
    }

    ingress {
        protocol = "tcp"
        from_port = 22
        to_port = 22
        security_groups = [aws_security_group.bastion.id]
    }

    egress {
        protocol = "-1"
        from_port = 0
        to_port = 0
        cidr_blocks = ["0.0.0.0/0"]
    }
}