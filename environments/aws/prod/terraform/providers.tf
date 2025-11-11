terraform {
    required_providers {
        aws = {
            source = "hashicorp/aws"
            version = "~> 5.92"
        }
    }

    backend "s3" {
      key = "dev/terraform.tfstate"
      encrypt = true
    }

    required_version = ">=1.2"
}

provider "aws" {
   region = var.region
   profile = var.env
}
