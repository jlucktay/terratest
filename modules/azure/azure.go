// Package azure allows interaction with resources on Microsoft Azure.
package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/pkg/errors"
)

// AzureSession is an object representing session for subscription
type AzureSession struct {
	SubscriptionID string
	Authorizer     autorest.Authorizer
}

func readJSON(path string) (*map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, errors.Wrap(err, "Can't open the file")
	}

	contents := make(map[string]interface{})
	err = json.Unmarshal(data, &contents)

	if err != nil {
		err = errors.Wrap(err, "Can't unmarshal file")
	}

	return &contents, err
}

func newSessionFromFile() (*AzureSession, error) {
	authorizer, err := auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)

	if err != nil {
		return nil, errors.Wrap(err, "Can't initialize authorizer")
	}

	authInfo, err := readJSON(os.Getenv("AZURE_AUTH_LOCATION"))

	if err != nil {
		return nil, errors.Wrap(err, "Can't get authinfo")
	}

	sess := AzureSession{
		SubscriptionID: (*authInfo)["subscriptionId"].(string),
		Authorizer:     authorizer,
	}

	return &sess, nil
}

func getGroups(sess *AzureSession) ([]string, error) {
	tab := make([]string, 0)
	var err error

	grClient := resources.NewGroupsClient(sess.SubscriptionID)
	grClient.Authorizer = sess.Authorizer

	for list, err := grClient.ListComplete(context.Background(), "", nil); list.NotDone(); err = list.Next() {
		if err != nil {
			return nil, errors.Wrap(err, "error traverising RG list")
		}
		rgName := *list.Value().Name
		tab = append(tab, rgName)
	}
	return tab, err
}

func getVM(sess *AzureSession, rg string, wg *sync.WaitGroup) {
	defer wg.Done()

	vmClient := compute.NewVirtualMachinesClient(sess.SubscriptionID)
	vmClient.Authorizer = sess.Authorizer

	for vm, err := vmClient.ListComplete(context.Background(), rg); vm.NotDone(); err = vm.Next() {
		if err != nil {
			log.Print("got error while traverising RG list: ", err)
		}

		i := vm.Value()
		tags := []string{}
		for k, v := range i.Tags {
			tags = append(tags, fmt.Sprintf("%s?%s", k, *v))
		}
		tagsS := strings.Join(tags, "%")

		if len(i.Tags) > 0 {
			fmt.Printf("%s,%s,%s,<%s>\n", rg, *i.Name, *i.ID, tagsS)
		} else {
			fmt.Printf("%s,%s,%s\n", rg, *i.Name, *i.ID)
		}
	}
}

func doStuff() {
	var wg sync.WaitGroup
	sess, err := newSessionFromFile()

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	groups, err := getGroups(sess)

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	wg.Add(len(groups))

	for _, group := range groups {
		go getVM(sess, group, &wg)
	}
	wg.Wait()
}

/*

 */

// NewGroupsClientE creates a Groups client.
func NewGroupsClientE(t *testing.T) (*resources.GroupsClient, error) {
	sess, err := newSessionFromFile()
	if err != nil {
		return nil, err
	}

	groupsClient := resources.NewGroupsClient(sess.SubscriptionID)
	groupsClient.Authorizer = sess.Authorizer
	groupsClient.AddToUserAgent("terratest/azure")

	return &groupsClient, nil
}

// GetTagsForResourceGroup returns all the tags for the given Resource Group.
func GetTagsForResourceGroup(t *testing.T, rgID string) map[string]string {
	tags, err := GetTagsForResourceGroupE(t, rgID)
	if err != nil {
		t.Fatal(err)
	}
	return tags
}

// GetTagsForResourceGroupE returns all the tags for the given Resource Group.
func GetTagsForResourceGroupE(t *testing.T, rgID string) (map[string]string, error) {
	client, err := NewGroupsClientE(t)
	if err != nil {
		return nil, err
	}

	_ = client

	return map[string]string{}, nil
}
