# Terraform provider for communicating with Provisioner API

- a fully-functional showcase for a Terraform provider written in Golang
- anonymized with names
- can be adapted and connected to your own platform API

Provider Usage
================
 - this provider is used with conjunction of terraform modules that managed virtual machine instances on all cloud providers used in Acme company
    - platform manages and cleans any specific metadata from 3rd-party services like DNS, Chef, Foreman, ...
 - Provisioner is a fictional hybrid platform-as-a-service managing on-premise and cloud virtual machines and resources
    - this providers allows cleanup of those related resources while still managing VMs on cloud with Terraform

## Provider initialization in Terraform
* initialize provider with:
  * api_endpoint or export PROVISIONER_API_ENDPOINT (optional since default is set)
  * api_key or export PROVISIONER_API_KEY, *export is strongly recommended* since it's a secret key!

## Provider build and test


```shell
go build -o terraform-provider-provisioner
```

```shell
make test
```

- `binary_output` directory contains the latest built provider version ready for usage (on current set `OS_ARCH` architecture) 

## Provider release

First, build the provider.

```shell
make release
```

- `bin/binary_output` directory contains the latest built provider version ready for usage (on current set `OS_ARCH` architecture)
- .zip file of the provider will be necessary for further steps!

## More

Check makefile for any other usage.