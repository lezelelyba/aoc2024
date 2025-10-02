# terraform/main.tf

terraform {
  required_providers {
    libvirt = {
      source  = "dmacvicar/libvirt"
      version = ">= 0.8.3"
    }
  }

  backend "local" {
    path = "/var/tmp/advent2024-terraform-kvm.tfstate"
  }
}

provider "libvirt" {
  uri = "qemu:///system"
}

resource "libvirt_network" "vm_bridge" {
  name      = "adv24-vm-br"
  bridge    = "adv24-vm-br"
  mode      = "nat"
  addresses = ["192.168.130.0/24"]
  dhcp {
    enabled = true
  }
  autostart = false
}

variable "vm_base_image_path" {
  # default = "/var/lib/libvirt/images/ubuntu-24.04-minimal-cloudimg-amd64.img"
  default = "/tmp/ubuntu-24.04-minimal-cloudimg-amd64.img"
}

variable "public_key_path" {
  default = "~/.ssh/id_ed25519.pub"
}
variable "private_key_path" {
  default = "~/.ssh/id_ed25519"
}

variable "vm_web_count" {
  default = 2
}

variable "vm_lb_mac" {
  default = "52:54:00:aa:00:01"
}

resource "libvirt_volume" "vm_base_volume" {
  name   = "vm-base.qcow2"
  pool   = "default"
  source = var.vm_base_image_path
  format = "qcow2"
}

resource "libvirt_volume" "vm_web_volume" {
  count    = var.vm_web_count
  name     = "vm_web_${count.index}.qcow2"
  pool     = "default"
  base_volume_id = libvirt_volume.vm_base_volume.id
  format   = "qcow2"
}
resource "libvirt_volume" "vm_lb_volume" {
  name     = "vm_lb.qcow2"
  pool     = "default"
  base_volume_id = libvirt_volume.vm_base_volume.id
  format   = "qcow2"
}

data "template_file" "vm_web_cloudinit_userdata" {
    template = "${file("../kvm/vm-web.cloudinit.userdata.template.yaml")}"
    vars = {
        hostname = "advent2024-web-docker"
        username = "ubuntu"
        password = "ubuntu"
        sshkey = file(var.public_key_path)
    }
}
data "template_file" "vm_web_cloudinit_network" {
    count = var.vm_web_count
    template = "${file("../kvm/vm-web.cloudinit.network.template.yaml")}"
    vars = {
        eth0_mac = "52:54:00:aa:bb:0${count.index}"
    }
}
data "template_file" "vm_lb_cloudinit_userdata" {
    template = "${file("../kvm/vm-lb.cloudinit.userdata.template.yaml")}"
    vars = {
        hostname = "advent2024-lb-docker"
        username = "caddy"
        password = "caddy"
        sshkey = file(var.public_key_path)
    }
}
data "template_file" "vm_lb_cloudinit_network" {
    template = "${file("../kvm/vm-lb.cloudinit.network.template.yaml")}"
    vars = {
        eth0_mac = var.vm_lb_mac
    }
}

resource "libvirt_cloudinit_disk" "vm_web_cloudinit" {
  count = var.vm_web_count
  name      = "vm-web.cloudinit.${count.index}.iso"
  # user_data = file("../kvm/vm-web.cloudinit.yaml")
  user_data = data.template_file.vm_web_cloudinit_userdata.rendered
  network_config = data.template_file.vm_web_cloudinit_network[count.index].rendered
}
resource "libvirt_cloudinit_disk" "vm_lb_cloudinit" {
  name      = "vm-lb.cloudinit.iso"
  user_data = data.template_file.vm_lb_cloudinit_userdata.rendered
  network_config = data.template_file.vm_lb_cloudinit_network.rendered
}

resource "libvirt_domain" "vm_web_docker" {
  count = var.vm_web_count

  name   = "vm-web-docker-${count.index}"
  memory = 4096 
  vcpu   = 2

  network_interface {
    network_id = libvirt_network.vm_bridge.id
    mac = "52:54:00:aa:bb:0${count.index}"
    wait_for_lease = true
  }

  disk {
    volume_id = libvirt_volume.vm_web_volume[count.index].id
  }

  # Unable to reference count.index in the depends on, but the bottom line should be enough
  # 
  # depends_on = [
  #   libvirt_cloudinit_disk.vm_web_cloudinit[count.index]
  # ]
  
  cloudinit = libvirt_cloudinit_disk.vm_web_cloudinit[count.index].id

  console {
    type        = "pty"
    target_port = "0"
  }

  graphics {
    type          = "vnc"
    listen_type   = "address"
    listen_address = "127.0.0.1"
  }

  cpu {
    mode = "host-model"
  }
}
resource "libvirt_domain" "vm_lb_docker" {
  name   = "vm-lb-docker"
  memory = 4096 
  vcpu   = 2

  network_interface {
    network_id = libvirt_network.vm_bridge.id
    mac = var.vm_lb_mac
    wait_for_lease = true
  }

  disk {
    volume_id = libvirt_volume.vm_lb_volume.id
  }
  
  cloudinit = libvirt_cloudinit_disk.vm_lb_cloudinit.id

  console {
    type        = "pty"
    target_port = "0"
  }

  graphics {
    type          = "vnc"
    listen_type   = "address"
    listen_address = "127.0.0.1"
  }

  cpu {
    mode = "host-model"
  }
}
