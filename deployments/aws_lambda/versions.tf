terraform {
  backend "s3" {}

  required_providers {
    archive = {
      source  = "hashicorp/archive"
      version = "2.2.0"
    }

    aws = {
      source  = "hashicorp/aws"
      version = "3.57.0"
    }
  }
}
