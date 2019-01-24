# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# See test/terraform_azure_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "example" {
  name     = "${var.rg_name}"
  location = "West Europe"

  tags {
    Name = "${var.rg_name}"
  }
}
