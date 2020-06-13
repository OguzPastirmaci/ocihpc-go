package main

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/example/helpers"
	"github.com/oracle/oci-go-sdk/resourcemanager"
	"github.com/jeffail/gabs"
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

jsonParsed, err := gabs.Pars




tf := resourcemanager.GetJobLogsRequest{
	JobId:                         applyJobResp.Id,
	TimestampGreaterThanOrEqualTo: &common.SDKTime{time.Now().Add(time.Second * -300)},
	SortOrder:                     "ASC",
}

resp, err := client.GetJobLogs(ctx, tf)
helpers.FatalIfError(err)

out, err := json.Marshal(resp.Items)
if err != nil {
	panic(err)
}
/*
	jsonParsed, err := gabs.ParseJSON(out)
	logs := jsonParsed.Path("data").Children()
	for _, log := range logs {
		msg := log.S("message").Data().(string)
		fmt.Println(msg)
	}
*/
//jsonParsed, err := gabs.ParseJSON([]byte(string(out)))

jsonParsed, err := gabs.ParseJSON(out)
logs := jsonParsed.Path("[]").Children()
for _, log := range logs {
	msg := log.S("type").Data().(string)
	fmt.Println(msg)
}

//fmt.Println(string(out))

var output string
output = jsonParsed.Path("type").Data().(string)

fmt.Println(output)