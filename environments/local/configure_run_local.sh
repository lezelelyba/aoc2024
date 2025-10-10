#!/bin/bash

# test1

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

DOCKER_IMAGE_TAG="advent2024.web:latest"
DOCKER_TAR_IMAGE="advent2024.web.latest.tar"
TMP_PREFIX="/tmp"

if [[ $1 == "cleanup" ]]; then
    rm -v $TMP_PREFIX/$DOCKER_TAR_IMAGE
    exit 0
fi

cd $SCRIPT_DIR

docker save -o "$TMP_PREFIX/$DOCKER_TAR_IMAGE" $DOCKER_IMAGE_TAG

ansible-playbook -i ansible/inventory.ini ansible/deploy.web.yml --extra-vars "tar_image=$TMP_PREFIX/$DOCKER_TAR_IMAGE image_tag=$DOCKER_IMAGE_TAG"
ansible-playbook -i ansible/inventory.ini ansible/deploy.lb.yml

rm $TMP_PREFIX/$DOCKER_TAR_IMAGE