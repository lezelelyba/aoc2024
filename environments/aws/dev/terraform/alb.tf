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

resource "aws_alb_listener" "front_end" {
    load_balancer_arn = aws_alb.main.arn 
    port = var.lb_tcp_port
    protocol = "HTTP"          
    default_action {
        target_group_arn = aws_alb_target_group.app.id
        type = "forward"
    }
}
