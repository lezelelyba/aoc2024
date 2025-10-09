output "alb_dns_name" {
  description = "Public URL of the Application Load Balancer"
  value       = "http://${aws_alb.main.dns_name}"
}

output "bastion_ip" {
  description = "Public IP of the Bastion server"
  value       = var.bastion ? "${aws_instance.bastion[0].public_ip}" : ""
}
output "private_host" {
  description = "Test Host ssh string"
  value       = var.private_host ? "ssh -i ${var.sshprivkeypath} -o ProxyCommand=\"ssh -i ${var.sshprivkeypath} -W %h:%p ec2-user@${aws_instance.bastion[0].public_ip}\" ec2-user@${aws_instance.private_host[0].private_ip}" : ""
}

output "ecs_cluster_name" {
  description = "ECS Cluster Name"
  value = aws_ecs_cluster.ecs_cluster.name
}

output "ecs_service_name" {
  description = "ECS Service Name"
  value = aws_ecs_service.app.name
}

resource "aws_ssm_parameter" "cd_config" {
  name  = var.ssm_path
  type  = "String"
  value = jsonencode({
    cluster = aws_ecs_cluster.ecs_cluster.name
    service = aws_ecs_service.app.name
  })
}