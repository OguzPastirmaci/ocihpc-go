package main

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/example/helpers"
	"github.com/oracle/oci-go-sdk/resourcemanager"
)

func main() {
	provider := common.DefaultConfigProvider()
	client, err := resourcemanager.NewResourceManagerClientWithConfigurationProvider(provider)
	helpers.FatalIfError(err)

	ctx := context.Background()

	jobLifecycle := resourcemanager.GetJobRequest{
		JobId: "ocid1.ormjob.oc1.iad.aaaaaaaavdlngdsgd6ncazajp6pkponyei5vnz5gkr4ijhisijme3fgb3rlq",
	}

	readResp, err := client.GetJob(ctx, jobLifecycle)

	fmt.Println(readResp.LifecycleState)

	tf := resourcemanager.GetJobTfStateRequest{
		JobId: "ocid1.ormjob.oc1.iad.aaaaaaaavdlngdsgd6ncazajp6pkponyei5vnz5gkr4ijhisijme3fgb3rlq",
	}

	readResp2, err := client.GetJobTfState(ctx, tf)

	fmt.Println(readResp2.Content)

}
