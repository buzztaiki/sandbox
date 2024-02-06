package main

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/buzztaiki/sandbox/azure/azure_sdk_go_list_user_and_groups/util"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
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
	userIterator, err := msgraphcore.NewPageIterator[models.Userable](
		usersResult, graphClient.RequestAdapter,
		models.CreateUserCollectionResponseFromDiscriminatorValue)
	if err != nil {
		return err
	}
	if err := userIterator.Iterate(ctx, func(user models.Userable) bool {
		fmt.Println(util.FromPtr(user.GetDisplayName()))
		return true
	}); err != nil {
		return err
	}

	return nil
}

func main() {
	run()
}
