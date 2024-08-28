terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
}

# Create an ECR repository
resource "aws_ecr_repository" "repo" {
  name = var.service_name
}

data "aws_ecr_authorization_token" "ecr_auth" {
}

# Tag the provided docker image with the ECR repository url
resource "docker_tag" "tag" {
  source_image = var.image
  target_image = aws_ecr_repository.repo.repository_url
}

# Push the tagged image to the ECR repository
resource "docker_registry_image" "push" {
  name = aws_ecr_repository.repo.repository_url
  triggers = {
    source_image_id = docker_tag.tag.source_image_id
  }
}

# Create a role for the lambda function
resource "aws_iam_role" "role" {
  name = var.service_name
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "basic-execution" {
  role       = aws_iam_role.role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Create a lambda function using the pushed image
resource "aws_lambda_function" "function" {
  function_name = var.service_name
  role          = aws_iam_role.role.arn
  image_uri     = "${aws_ecr_repository.repo.repository_url}@${docker_registry_image.push.sha256_digest}"
  package_type  = "Image"
  timeout = 30
  environment {
    variables = var.environment
  }

  depends_on = [docker_registry_image.push]
}
