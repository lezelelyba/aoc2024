Demo project to showcase IaC. "Backend" is solver for [Advent of Code 2024](https://adventofcode.com/)

Utilizes:
  - GitHub Actions
    - For Docker container image creation
    - For AWS Deployment utilizing Terraform
  - Terraform
    - For local KVM deployment
    - For AWS bootstrap
    - For AWS deployment
  - Ansible
    - For running container in local KVM deployment

Requirements:
  - Terraform installed locally for bootstrap
  - AWS account
  - AWS CLI installed
  - Defined AWS user which can boostrap the environment

Manual Tasks:
  - Create User in AWS and configure AWS CLI
  - Modify the AWS CLI config to contain environment specified in aws_bootstrap/variables.tf ${env}
  - Clone repo
  - Bootstrap the environment
    - creates Github Action for dev and master branch ./github/workflows/deploy.yml
  - Push to repo
  - Run Deploy Github action to build the AWS environment