terraform {
  backend "s3" {
    key = "stage/terraform.tfstate"
    encrypt = true
  }
}