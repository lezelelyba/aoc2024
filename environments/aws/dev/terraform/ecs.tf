resource "aws_ecs_cluster" "ecs_cluster" {
   name = "${var.env}-aoc2024-ecs-cluster" 
}

resource "aws_iam_role" "ecs_task_execution_role" {
    name = "${var.env}Aoc2024ExecutionRole"

    assume_role_policy = jsonencode({
        Version = "2012-10-17"
        Statement = [{
            Action = "sts:AssumeRole"
            Effect = "Allow"
            Principal = {
                Service = "ecs-tasks.amazonaws.com"
            }
        }]
    })
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_policy" {
    role = aws_iam_role.ecs_task_execution_role.name
    policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_ecs_task_definition" "ecs_task" {
    family = "${var.env}-aoc2024-ecs-task"
    network_mode = "awsvpc"

    requires_compatibilities = ["FARGATE"]

    cpu = 256
    memory = 512

    execution_role_arn = aws_iam_role.ecs_task_execution_role.arn

    container_definitions = jsonencode([
        {
            name = "aoc2024"
            image = var.docker_image
            essential = true
            portMappings = [
                {
                    containerPort = tonumber(var.app_tcp_port)
                    hostPort = tonumber(var.app_tcp_port)
                    protocol = "tcp"
                }
            ]
        }
    ])
}
resource "aws_ecs_service" "app" {
    name = "${var.env}-aoc2024-app"
    cluster = aws_ecs_cluster.ecs_cluster.id 
    task_definition = aws_ecs_task_definition.ecs_task.arn
    desired_count = 2
    launch_type = "FARGATE"

    network_configuration {
        subnets = aws_subnet.private[*].id
        security_groups = [aws_security_group.ecs.id]
    }

    load_balancer {
        target_group_arn = aws_alb_target_group.app.arn
        container_name = "aoc2024" 
        container_port = var.app_tcp_port
    }

    depends_on = [aws_alb.main]
}
