package main

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/buzztaiki/sandbox/azure/azure_sdk_go_list_user_and_groups/util"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

func run() error {
	ctx := context.Background()

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}
	graphClient, err := msgraphsdkgo.NewGraphServiceClientWithCredentials(cred, nil)
	if err != nil {
		return err
	}

	usersResult, err := graphClient.Users().Get(ctx, nil)
	if err != nil {
		return err
	}
	for {
		for _, user := range usersResult.GetValue() {
			fmt.Println(util.FromPtr(user.GetDisplayName()))
		}

		nextLink := usersResult.GetOdataNextLink()
		if nextLink == nil {
			break
		}
		usersResult, err = graphClient.Users().WithUrl(*nextLink).Get(ctx, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	run()
}
