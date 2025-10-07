resource "aws_security_group" "lb" {
    vpc_id = aws_vpc.vpc.id
    name = "${var.env}-aoc2024-lb-sg"

    tags = {
        name = "${var.env}-aoc2024-lb-sg"
        environment = var.env
    }

    ingress {
        protocol = "tcp"
        from_port = tonumber(var.lb_tcp_port)
        to_port = tonumber(var.lb_tcp_port)
        cidr_blocks = ["0.0.0.0/0"]
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