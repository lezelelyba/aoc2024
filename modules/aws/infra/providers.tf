terraform {
    required_providers {
        aws = {
            source = "hashicorp/aws"
            version = "~> 5.92"
        }
    }

    required_version = ">=1.2"
}

provider "aws" {
   alias = "prov1"
   region = var.region
   profile = var.env
}