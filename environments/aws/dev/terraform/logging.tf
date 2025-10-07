resource "aws_cloudwatch_log_group" "group" {
  name = "${var.env}-aoc2024-log-group"
}