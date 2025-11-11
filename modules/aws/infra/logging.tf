// cloud watch group for the ecs containers
resource "aws_cloudwatch_log_group" "ecs_log_group" {
  name = "${var.env}/ecs/aoc2024"
  retention_in_days = 14
}