package azure

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/gobuffalo/envy"
)

func getGroupsClientWithAuthFile() resources.GroupsClient {
	subscriptionID, errSub := envy.MustGet("AZURE_SUBSCRIPTION_ID")
	if errSub != nil {
		log.Fatalf("couldn't find subscription in environment: %v\n", errSub)
	}
	groupsClient := resources.NewGroupsClient(subscriptionID)
	// requires env var AZURE_AUTH_LOCATION set to output of `az ad sp create-for-rbac --sdk-auth`
	a, errAuth := auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	if errAuth != nil {
		log.Fatalf("failed to initialize authorizer: %v\n", errAuth)
	}
	groupsClient.Authorizer = a
	groupsClient.AddToUserAgent("terratest/azure")
	return groupsClient
}

// https://docs.microsoft.com/en-au/go/azure/
// https://docs.microsoft.com/en-au/go/azure/azure-sdk-go-authorization
// https://docs.microsoft.com/en-au/go/azure/azure-sdk-go-qs-vm

// https://github.com/Azure-Samples/azure-sdk-for-go-samples/blob/master/resources/groups.go#L26
