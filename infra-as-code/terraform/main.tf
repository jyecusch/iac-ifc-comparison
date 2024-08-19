terraform {
  required_version = ">= 1.0.0"
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

data "aws_ecr_authorization_token" "token" {
}

provider "docker" {
  registry_auth {
    address  = data.aws_ecr_authorization_token.token.proxy_endpoint
    username = data.aws_ecr_authorization_token.token.user_name
    password = data.aws_ecr_authorization_token.token.password
  }

}

module "subscriber" {
  source       = "./modules/service"
  service_name = "iac-subscriber"
  image        = "subscriber:latest"
  environment  = {}
}

locals {
  bus_source_id = "custom.source"
}

module "event_bus" {
  source    = "./modules/eventbridge-pubsub"
  bus_name  = "iac-bus"
  source_id = local.bus_source_id
  lambda_subscribers = {
    subfunc = module.subscriber.function_arn
  }
}

module "publisher" {
  source       = "./modules/service"
  service_name = "iac-publisher"
  image        = "publisher:latest"
  environment = {
    # If these env var names don't exactly match the expected env var name in the publisher code
    #   the publisher will not be able to publish to the EventBridge bus.
    EVENT_BUS_NAME  = module.event_bus.bus_name
    EVENT_SOURCE_ID = local.bus_source_id
  }
}

resource "aws_iam_role_policy" "policy" {
  # The Terraform code has no idea if the Lambda function still needs these permissions.
  #   It's up to the Terraform developer to have a detailed understanding of the code's requirements,
  #   then ensure they're reflected accurately in this deployment code,
  #   including removing these permissions if they are no longer needed.
  role = module.publisher.role_name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = {
      Effect   = "Allow"
      Action   = "events:PutEvents"
      Resource = module.event_bus.bus_arn
    }
  })
}

module "api" {
  source = "./modules/api"
  name   = "my-api"
  spec   = templatefile("../api_spec.json", { trigger_uri = module.publisher.invoke_arn })
  target_lambda_functions = {
    "publisher" = module.publisher.function_name
  }
}

output "api_url" {
  value = module.api.endpoint
}
