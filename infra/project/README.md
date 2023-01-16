# Project

## Overview

This directory contains the main bootstrap IaC for the project using Terraform.

This includes:
- Google Cloud APIs to enable
- Storage bucket for Terraform state
- Storage bucket for project data
- Secrets captured in Secret Manager for use by other services
	- ACLOUDGURU_API_KEY
	- ACLOUDGURU_CONSUMER_ID
	- SLACK_BOT_TOKEN
	- SLACK_SIGNING_SECRET
