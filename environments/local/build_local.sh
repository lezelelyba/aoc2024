#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

QEMU_PREFIX="/var/lib/libvirt/images"
CLOUD_IMG="ubuntu-24.04-minimal-cloudimg-amd64.img"

DOCKER_IMAGE_TAG="advent2024.web:latest"
DOCKER_TAR_IMAGE="advent2024.web.latest.tar"
TMP_PREFIX="/tmp"

cd $SCRIPT_DIR/terraform

if [[ $1 == "destroy" ]]; then
    terraform destroy
    rm $TMP_PREFIX/$CLOUD_IMG
    exit 0
fi

cd $TMP_PREFIX

if [[ ! -f $CLOUD_IMG ]]; then
    wget -O $CLOUD_IMG https://cloud-images.ubuntu.com/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64.img
fi

cd $SCRIPT_DIR/terraform

terraform init
terraform plan
terraform apply