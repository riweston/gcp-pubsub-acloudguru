# Application

## Overview

This directory contains the application IaC for this service deployed using Terraform.

To reduce repeated code and to promote reusability most of the logic is wrapped up in the
`infra/application/_modules/cloudfunction` module. This module has a couple of different patterns of
deployment depending on the parameters passed to it.

These include:
- Public HTTP trigger
- Private HTTP trigger triggered by Pub/Sub topic subscription

Implicitly this module also creates the Pub/Sub topic and subscription if required.

Important parameters to note are:
- `pubsub_topic` - The name of the Pub/Sub topic to create
- `pubsub_subscription` - The name of the Pub/Sub subscription to create
- `secret_environment_variables` - A list of environment variables to be populated from Secret Manager (note the name of the secret is the same as the environment variable name), this includes assignment of IAM permissions for the function service account to access the secrets
