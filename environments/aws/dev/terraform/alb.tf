resource "aws_alb" "main" {
    
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

    tags = {
        name = "${var.env}-aoc2024-lb-grp"
        environment = var.env
    }
}

resource "aws_alb_listener" "front_end" {
   load_balancer_arn = aws_alb.main.arn 
   port = var.app_tcp_port
   protocol = "HTTP"

   default_action {
    target_group_arn = aws_alb_target_group.app.id
    type = "forward"
   }
}
