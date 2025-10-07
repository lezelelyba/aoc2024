output "alb_dns_name" {
  description = "Public URL of the Application Load Balancer"
  value       = "http://${aws_alb.main.dns_name}"
}
output "bastion_ip" {
  description = "Public IP of the Bastion server"
  value       = "${aws_instance.bastion.public_ip}"
}
output "test_host" {
  description = "Test Host ssh string"
  value       = "ssh -i ${var.sshprivkeypath} -o ProxyCommand=\"ssh -i ${var.sshprivkeypath} -W %h:%p ec2-user@${aws_instance.bastion.public_ip}\" ec2-user@${aws_instance.test_host.private_ip}"
}