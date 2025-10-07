output "alb_dns_name" {
  description = "Public URL of the Application Load Balancer"
  value       = "http://${aws_alb.main.dns_name}"
}