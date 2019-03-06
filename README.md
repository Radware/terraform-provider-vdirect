<img src="https://www.radware.com/RadwareSite/MediaLibraries/Images/logo.svg" width="300px">

# Overview
A [Terraform](terraform.io) provider for [Radware vDirect](https://www.radware.com/products/vdirect/).

# Requirements
-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

# Radware vDirect requirements
- This provider uses the vDirect REST API, make sure that it is installed in your environment.
- All the resources are validated with vDirect 4.6

# Using the Provider

Place the plugin in your plugins directory, run terraform init to initialize it.
