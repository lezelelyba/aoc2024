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
  - Terraform
    - Install locally for bootstrap and AWS environment creation
  - AWS
    - create account
    - Define user who can bootstrap the environment
    - Install CLI
      - Modify the AWS CLI config to contain environment specified in <code>aws_bootstrap/variables.tf ${envs}</code>
      - Specify the User
    - Register domain with AWS
  - GitHub
    - specify AWS Secrets for GH Actions
    - Create OAuth
    - Obtain client id and client secret for OAuth

Tasks:
  make ENVIRONMENT=dev apply TF_APPLYPARAMS="-var-file=\"ecs.tfvars\"
  make ENVIRONMENT=stage apply TF_APPLYPARAMS="-var-file=\"ecs.tfvars\"
  make ENVIRONMENT=prod apply TF_APPLYPARAMS="-var-file=\"ecs.tfvars\"