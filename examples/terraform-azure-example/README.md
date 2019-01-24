# Terraform Azure Example

This folder contains a simple Terraform module that deploys resources in [Azure](https://azure.microsoft.com) to
demonstrate how you can use Terratest to write automated tests for your Azure Terraform code. This module deploys a
[Resource Group](https://docs.microsoft.com/azure/azure-resource-manager/resource-group-overview#resource-groups) and
gives that RG a `Name` tag with the value specified in the `rg_name` variable.

Check out [test/terraform_azure_example_test.go](/test/terraform_azure_example_test.go) to see how you can write
automated tests for this module.

Note that the resources deployed in this module doesn't actually do anything; RGs are mainly abstract containers for
holding other resources, for organisation and compliance purposes. For slightly more complicated, real-world examples of
Terraform modules, see [terraform-http-example](/examples/terraform-http-example) and
[terraform-ssh-example](/examples/terraform-ssh-example).

**WARNING**: This module and the automated tests for it deploy real resources into your Azure account which can cost you
money. The resources are all part of [Azure's always free tier](https://azure.microsoft.com/free/), so if you haven't
used that up, it should be free, but you are completely responsible for all Azure charges.

## Running this module manually

1. Sign up for [Azure](https://azure.microsoft.com).
1. Configure your Azure credentials using one of the [supported methods for Azure CLI
  tools](https://docs.microsoft.com/cli/azure/authenticate-azure-cli).
1. Set the active Azure subscription with `az account set --subscription <guid>` so that test resources are deployed
  into the correct subscription.
1. Set the default region in Azure you want to use with `az configure`.
1. Install [Terraform](https://www.terraform.io) and make sure it's on your `PATH`.
1. Run `terraform init`.
1. Run `terraform apply`.
1. When you're done, run `terraform destroy`.

## Running automated tests against this module

1. Sign up for [Azure](https://azure.microsoft.com).
1. Configure your Azure credentials using one of the [supported methods for Azure CLI
  tools](https://docs.microsoft.com/cli/azure/authenticate-azure-cli).
1. Set the active Azure subscription with `az account set --subscription <guid>` so that test resources are deployed
  into the correct subscription.
1. Install [Terraform](https://www.terraform.io) and make sure it's on your `PATH`.
1. Install [Golang](https://golang.org) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `dep ensure`
1. `go test -v -run TestTerraformAzureExample`
