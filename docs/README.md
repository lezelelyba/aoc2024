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

Tasks:
  - Terraform
    - Install locally for bootstrap and AWS environment creation
  - AWS
    - Create account
    - Define user who can bootstrap the environment
    - Install CLI
      - Modify the AWS CLI config to contain environment specified in <code>aws_bootstrap/variables.tf ${envs}</code>
      - Specify the User
    - Register domain with AWS
  - if OAuth is desired
  - GitHub
    - Specify AWS Secrets for GH Actions
    - Create OAuth
    - Obtain client id and client secret for OAuth
  - Azure
    - Install Azure CLI
    - az role assignment create --assignee {user-id} --role "Storage Blob Data Contributor" --scope /subscriptions/{subscription_id}
      - to be able to push the image to ACR - before CI is set
    - az role assignment create --assignee {user-id} --role "Key Vault Secrets Officer" --scope /subscriptions/{subscription_id}
      - to be able to push secrets to vault

Commands:
  - make ENRIRONMENT=bootstrap init
  - make ENRIRONMENT=bootstrap apply
  - make ENVIRONMENT=dev init
  - make ENVIRONMENT=dev apply
  - OR if you want to specify options for the container
  - make ENVIRONMENT=dev apply TF_APPLYPARAMS="-var-file=\"ecs.tfvars\""