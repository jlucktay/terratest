package azure

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/classic/management"
	"github.com/Azure/azure-sdk-for-go/services/classic/management/location"
	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
)

// You can set this environment variable to force Terratest to use a specific region rather than a random one. This is
// convenient when iterating locally.
const regionOverrideEnvVarName = "TERRATEST_REGION"

var stableRegions = []string{
	"australiacentral",
	"australiacentral2",
	"australiaeast",
	"australiasoutheast",
	"brazilsouth",
	"canadacentral",
	"canadaeast",
	"centralindia",
	"centralus",
	"eastasia",
	"eastus",
	"eastus2",
	"francecentral",
	"francesouth",
	"japaneast",
	"japanwest",
	"koreacentral",
	"koreasouth",
	"northcentralus",
	"northeurope",
	"southcentralus",
	"southeastasia",
	"southindia",
	"uksouth",
	"ukwest",
	"westcentralus",
	"westeurope",
	"westindia",
	"westus",
	"westus2",
}

// GetRandomRegion gets a randomly chosen Azure region. If approvedRegions is not empty, this will be a region from the
// approvedRegions list; otherwise, this method will fetch the latest list of regions from the Azure APIs and pick one
// of those. If forbiddenRegions is not empty, this method will make sure the returned region is not in the
// forbiddenRegions list.
func GetRandomRegion(t *testing.T, approvedRegions []string, forbiddenRegions []string) string {
	region, err := GetRandomRegionE(t, approvedRegions, forbiddenRegions)
	if err != nil {
		t.Fatal(err)
	}
	return region
}

// GetRandomRegionE gets a randomly chosen Azure region. If approvedRegions is not empty, this will be a region from
// the approvedRegions list; otherwise, this method will fetch the latest list of regions from the Azure APIs and pick
// one of those. If forbiddenRegions is not empty, this method will make sure the returned region is not in the
// forbiddenRegions list.
func GetRandomRegionE(t *testing.T, approvedRegions []string, forbiddenRegions []string) (string, error) {
	regionFromEnvVar := os.Getenv(regionOverrideEnvVarName)
	if regionFromEnvVar != "" {
		logger.Logf(t, "Using Azure region %s from environment variable %s", regionFromEnvVar, regionOverrideEnvVarName)
		return regionFromEnvVar, nil
	}

	regionsToPickFrom := approvedRegions

	if len(regionsToPickFrom) == 0 {
		allRegions, err := GetAllAzureRegionsE(t)
		if err != nil {
			return "", err
		}
		regionsToPickFrom = allRegions
	}

	regionsToPickFrom = collections.ListSubtract(regionsToPickFrom, forbiddenRegions)
	region := random.RandomString(regionsToPickFrom)

	logger.Logf(t, "Using region %s", region)
	return region, nil
}

// GetAllAzureRegions gets the list of Azure regions available in this account.
func GetAllAzureRegions(t *testing.T) []string {
	out, err := GetAllAzureRegionsE(t)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetAllAzureRegionsE gets the list of Azure regions available in this account.
func GetAllAzureRegionsE(t *testing.T) ([]string, error) {
	logger.Log(t, "Looking up all Azure regions available")

	ac := management.NewAnonymousClient()
	lc := location.NewClient(ac)

	locs, err := lc.ListLocations()
	if err != nil {
		return nil, err
	}

	regions := []string{}
	for _, region := range locs.Locations {
		regions = append(regions, region.Name)
	}

	return regions, nil
}
