Demo project to showcase IaC. "Backend" is solver for [Advent of Code 2024](https://adventofcode.com/)

Utilizes:
  - GitHub Actions
    - For CI: Docker container image creation
    - For CD: Pushing updated docker image to ECS
  - Terraform
    - For local KVM deployment
    - For AWS bootstrap
    - For AWS deployment
  - Ansible
    - For running container in local KVM deployment

Requirements:
  - Terraform installed locally for bootstrap and AWS environment creation
  - AWS account
  - AWS CLI installed
  - Defined AWS user which can boostrap the environment

Manual Tasks:
  - Clone repo
  - Create User in AWS and configure AWS CLI
  - Modify the AWS CLI config to contain environment specified in <code>aws_bootstrap/variables.tf ${envs}</code>
  - Bootstrap the environment via <code>make bootstrap</code>
    - Creates:
      - S3 bucket for TF state
      - Dynamo DB for TF locks
      - GH Action role for CD
  - Create the environemtn via <code>make dev</code> or <code>make prod</code>