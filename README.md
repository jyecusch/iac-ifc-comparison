# Comparing Infrastructure Deployment

## Overview

This project contains two implementations of the same cloud application, the first uses common Infrastructure-as-Code (infra-as-code) techniques for deployment, the second uses Infrastructure-from-Code (infra-from-code) to highlight the differences.

For each implementation there are two versions of the application:
  - Version 1 uses AWS SNS for async messaging between AWS Lambda Functions
  - Version 2 substitutes AWS SNS for AWS EventBridge

Each version is represented by a single Git Commit, you can review the commits individually to see the implementations and the diff between v1 and v2.

Both examples are written in Go and deploy the following resources to the cloud:

### Version 1 (first commit)
  - 2x AWS Lambda Functions, deployed as Docker Containers
    - A `publisher` (pub) function, triggered via an API and publishes messages to an SNS topic
    - A `subscriber` (sub) function, triggered by a subscription to the SNS topic
  - An HTTP API Gateway, used to trigger the 'publisher' function
  - An SNS Topic to facilitate async messaging between the services

### Version 2 (second commit)
  - 2x AWS Lambda Functions, deployed as Docker Containers
    - A `publisher` (pub) function, triggered via an API and sends messages to an EventBridge event bus
    - A `subscriber` (sub) function, triggered by an event rule for messages from the event bus
  - An HTTP API Gateway, used to trigger the 'publisher' function
  - An EventBridge (CloudWatch) Event Bus to facilitate async messaging between the services

The purpose of these examples it to demonstrate the improved separation of concerns when using Infrastructure-from-Code (which is an enhancement of Infrastructure-as-Code, not a replacement)

> *Note:* Some areas of interest are highlighted inline using code comments

Since the deployment automation modules (e.g. Terraform Modules) need to pre-exist or be modified in either case, we've included all the required modules in the first commit (e.g. the EventBridge Terraform Module & Nitric EventBridge Provider Extension).

This makes the application layer changes clearer in the commit diff.

## Why swap SNS for EventBridge?

The services we're swapping don't really matter, it could be any two services. However, the replacement of a managed cloud service highlights how easy it is for a lack of separation of concerns to develop in applications using typical Infrastructure-as-Code techniques.

Swapping these similar services (at least the way we're using them here) requires broad changes in the infra-as-code example. Both functions change, their tests change and their deployment automation code changes. There is also the potential for common errors like a typo in an env var name to break the application - which is tricky to catch before it's deployed.

By comparison, the Infrastructure-from-Code example requires no changes at the application layer. This is because IfC builds on IaC, introducing a layer of separation, without sacrificing control. 

## Running the examples

Both examples include a `Makefile` with commands for testing and deployment.

### Infra-as-Code

You'll need the following dependencies:

  - [Docker](https://docs.docker.com/get-docker/)
  - [Terraform CLI](https://developer.hashicorp.com/terraform/install)
  - [Go v1.22](https://go.dev/doc/install)

Deploying to AWS

```bash
cd infra-as-code

# deploy the project
make terraform-init
make terraform-apply
```

### Infra-from-Code

You'll need the following dependencies:

  - [Docker](https://docs.docker.com/get-docker/)
  - [Nitric CLI](https://nitric.io/docs/getting-started/installation)
  - [Pulumi](https://www.pulumi.com/docs/cli/)
  - [Go v1.22](https://go.dev/doc/install)

> Nitric is capable of generating IaC with many tools, including both Terraform and Pulumi. This example uses Pulumi, but could equally have used Terraform without impacting the results.

Deploying to AWS

```bash
cd infra-from-code

# deploy the project
make up
```

The nitric project can also run locally

```bash
cd infra-from-code

# run the project locally
make run
```

> The standard Nitric providers use SNS in AWS by default, just like we created a Terraform Module for EventBridge, we've included a nitric AWS provider extension for EventBridge. This will be built before deployment.

The Nitric provider extension in `./eventbridge-provider` is a good example of _how_ nitric adds separation between cloud and application concerns.
