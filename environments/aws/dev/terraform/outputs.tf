output "alb_dns_name_aws" {
  description = "Public URL of the Application Load Balancer - AWS generated"
  value       = module.infra.alb_dns_name_aws
}
output "alb_dns_name" {
  description = "Public URL of the Application Load Balancer"
  value       = module.infra.alb_dns_name
}

output "bastion_ip" {
  description = "Public IP of the Bastion server"
  value       = module.infra.bastion_ip
}

output "private_host" {
  description = "Test Host ssh string"
  value       = module.infra.private_host
}