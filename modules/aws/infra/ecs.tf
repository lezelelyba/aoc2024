// allow ECS task to pull secrets from vault
// aim role
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

// aim role policy
resource "aws_iam_role_policy" "ecs_task_execution_secrets_policy" {
    name = "${var.env}Aoc2024ExecutionSecrets"
    role = aws_iam_role.ecs_task_execution_role.id

    policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
          {
            Effect = "Allow"
            Action = ["secretsmanager:GetSecretValue"]
            Resource = "*"
          }
        ]
    })
}

// aim role <> policy
resource "aws_iam_role_policy_attachment" "ecs_task_execution_policy" {
    role = aws_iam_role.ecs_task_execution_role.name
    policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

// create ECS cluster
resource "aws_ecs_cluster" "ecs_cluster" {
   name = "${var.env}-aoc2024-ecs-cluster" 
}

// define the ECS task
resource "aws_ecs_task_definition" "ecs_task" {
    family = "${var.env}-aoc2024-ecs-task"
    network_mode = "awsvpc"

    requires_compatibilities = ["FARGATE"]

    cpu = 256
    memory = 512

    // assing aim role
    execution_role_arn = aws_iam_role.ecs_task_execution_role.arn

    container_definitions = jsonencode([
        {
            name = "aoc2024"
            image = var.docker_image
            essential = true
            // port mapping
            portMappings = [
                {
                    containerPort = tonumber(var.app_tcp_port)
                    hostPort = tonumber(var.app_tcp_port)
                    protocol = "tcp"
                }
            ]

            // env variables
            environment = [
                for k, v in var.ecs_app_env_map : {
                    name = k
                    value = v
                }
            ]

            // env secrets 
            secrets = [
                for i in local.ecs_app_env_map_secret_keys : {
                    name = i
                    valueFrom = aws_secretsmanager_secret.container_env_secret["${var.env}-${i}"].arn
                }
            ]

            // cloud watch
            logConfiguration = {
                logDriver = "awslogs"
                options = {
                  awslogs-group = aws_cloudwatch_log_group.ecs_log_group.name
                  awslogs-region = var.region
                  awslogs-stream-prefix = "${var.env}-ecs-aoc2024"
                }
            }
        }
    ])
}

// define ECS service
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

// define secret keys
resource "aws_secretsmanager_secret" "container_env_secret" {
    // single ssm, prefix keys with env to avoid conflict
    for_each = { for k in local.ecs_app_env_map_secret_keys: "${var.env}-${k}" => k }

    name = each.key
    recovery_window_in_days = 0
}

// save secret values
resource "aws_secretsmanager_secret_version" "container_env_secret_value" {
    for_each = aws_secretsmanager_secret.container_env_secret

    secret_id = each.value.id
    secret_string = var.ecs_app_env_map_secret[replace(each.key, "${var.env}-", "")]
}