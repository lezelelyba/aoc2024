resource "local_file" "ansible_inventory" {
  filename = "../ansible/inventory.ini"

  content = <<EOT
[webservers]
${join("\n", [
  for i in range(length(libvirt_domain.vm_web_docker)):
    "${libvirt_domain.vm_web_docker[i].name} ansible_host=${libvirt_domain.vm_web_docker[i].network_interface[0].addresses[0]} ansible_user=ubuntu ansible_ssh_common_args='-o StrictHostKeyChecking=no'"
])}
[loadbalancers]
${libvirt_domain.vm_lb_docker.name} ansible_host=${libvirt_domain.vm_lb_docker.network_interface[0].addresses[0]} ansible_user=caddy ansible_ssh_common_args='-o StrictHostKeyChecking=no'
EOT

  depends_on = [libvirt_domain.vm_web_docker, libvirt_domain.vm_lb_docker]
}