package test

import (
	"fmt"
	"testing"

	"github.com/jlucktay/terratest/modules/azure"
	"github.com/jlucktay/terratest/modules/random"
	"github.com/jlucktay/terratest/modules/terraform"
)

// An example of how to test the Terraform module in examples/terraform-azure-example using Terratest.
func TestTerraformAzureExample(t *testing.T) {
	t.Parallel()

	// Give this RG a unique ID for a name tag so we can distinguish it from any other existing RG in your Azure
	// subscription.
	expectedName := fmt.Sprintf("terratest-azure-example-%s", random.UniqueId())

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located.
		TerraformDir: "../examples/terraform-azure-example",

		// Variables to pass to our Terraform code using -var options.
		Vars: map[string]interface{}{
			"rg_name": expectedName,
		},

		// Environment variables to set when running Terraform.
		EnvVars: map[string]string{},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable.
	resourceGroupID := terraform.Output(t, terraformOptions, "rg_id")

	_ = azure.GetTagsForResourceGroup(t, resourceGroupID)

	/*
		aws.AddTagsToResource(t, awsRegion, instanceID, map[string]string{"testing": "testing-tag-value"})

		// Look up the tags for the given Instance ID
		instanceTags := aws.GetTagsForEc2Instance(t, awsRegion, instanceID)

		testingTag, containsTestingTag := instanceTags["testing"]
		assert.True(t, containsTestingTag)
		assert.Equal(t, "testing-tag-value", testingTag)

		// Verify that our expected name tag is one of the tags
		nameTag, containsNameTag := instanceTags["Name"]
		assert.True(t, containsNameTag)
		assert.Equal(t, expectedName, nameTag)
	*/
}
