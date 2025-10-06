# Manually created for testing

terraform {
  backend "s3" {
    bucket         = "bucket"
    key            = "dev/terraform.tfstate"
    region         = "eu-central-1"
    dynamodb_table = "tf-lock-db"
    encrypt        = true
  }
}