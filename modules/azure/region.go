package azure

// import (
// 	"os"
// 	"testing"

// 	"github.com/jlucktay/terratest/modules/collections"
// 	"github.com/jlucktay/terratest/modules/logger"
// 	"github.com/jlucktay/terratest/modules/random"
// )

// // GetRandomRegion gets a randomly chosen Azure region. If approvedRegions is not empty, this will be a region from the
// // approvedRegions list; otherwise, this method will fetch the latest list of regions from the Azure APIs and pick one
// // of those. If forbiddenRegions is not empty, this method will make sure the returned region is not in the
// // forbiddenRegions list.
// func GetRandomRegion(t *testing.T, approvedRegions []string, forbiddenRegions []string) string {
// 	region, err := GetRandomRegionE(t, approvedRegions, forbiddenRegions)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	return region
// }

// // GetRandomRegionE gets a randomly chosen Azure region. If approvedRegions is not empty, this will be a region from
// // the approvedRegions list; otherwise, this method will fetch the latest list of regions from the Azure APIs and pick
// // one of those. If forbiddenRegions is not empty, this method will make sure the returned region is not in the
// // forbiddenRegions list.
// func GetRandomRegionE(t *testing.T, approvedRegions []string, forbiddenRegions []string) (string, error) {
// 	regionFromEnvVar := os.Getenv(regionOverrideEnvVarName)
// 	if regionFromEnvVar != "" {
// 		logger.Logf(t, "Using AWS region %s from environment variable %s", regionFromEnvVar, regionOverrideEnvVarName)
// 		return regionFromEnvVar, nil
// 	}

// 	regionsToPickFrom := approvedRegions

// 	if len(regionsToPickFrom) == 0 {
// 		allRegions, err := GetAllAwsRegionsE(t)
// 		if err != nil {
// 			return "", err
// 		}
// 		regionsToPickFrom = allRegions
// 	}

// 	regionsToPickFrom = collections.ListSubtract(regionsToPickFrom, forbiddenRegions)
// 	region := random.RandomString(regionsToPickFrom)

// 	logger.Logf(t, "Using region %s", region)
// 	return region, nil
// }
