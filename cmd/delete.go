/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/example/helpers"
	"github.com/oracle/oci-go-sdk/resourcemanager"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		region, _ := cmd.Flags().GetString("region")
		provider := common.DefaultConfigProvider()
		solution, _ := cmd.Flags().GetString("solution")
		client, err := resourcemanager.NewResourceManagerClientWithConfigurationProvider(provider)
		if err != nil {
			panic(err)
		}
		helpers.FatalIfError(err)

		ctx := context.Background()
		stackID := getStackInfo("StackID")

		createDestroyJob(ctx, provider, client, stackID, region, solution)

	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringP("region", "r", "", "The region to deploy to")
	deleteCmd.MarkFlagRequired("region")

	deleteCmd.Flags().StringP("solution", "s", "", "Solution to delete")
	deleteCmd.MarkFlagRequired("solution")
}

func deleteStack(ctx context.Context, stackID string, client resourcemanager.ResourceManagerClient, solution string) {

	req := resourcemanager.DeleteStackRequest{
		StackId: common.String(stackID),
	}

	_, err := client.DeleteStack(ctx, req)
	helpers.FatalIfError(err)

	fmt.Println("Stack deletion")
}

func createDestroyJob(ctx context.Context, provider common.ConfigurationProvider, client resourcemanager.ResourceManagerClient, stackID string, region string, solution string) string {

	destroyJobReq := resourcemanager.CreateJobRequest{
		CreateJobDetails: resourcemanager.CreateJobDetails{
			StackId:   common.String(stackID),
			Operation: "DESTROY",
			JobOperationDetails: resourcemanager.CreateDestroyJobOperationDetails{
				ExecutionPlanStrategy: "AUTO_APPROVED",
			},
		},
	}

	destroyJobResp, err := client.CreateJob(ctx, destroyJobReq)
	//fmt.Printf("Job ID of the destroy job is: %s\n", *destroyJobResp.Id)

	if err != nil {
		fmt.Println("Submission of destroy job failed", err)
		os.Exit(1)
	}

	jobLifecycle := resourcemanager.GetJobRequest{
		JobId: destroyJobResp.Id,
	}

	start := time.Now()

	for {
		elapsed := int(time.Since(start).Seconds())
		readResp, err := client.GetJob(ctx, jobLifecycle)

		if err != nil {
			fmt.Println("Destroy job failed", err)
			os.Exit(1)
		}

		fmt.Printf("Deleting solution: %s [%dmin %dsec]\n", solution, elapsed/60, elapsed%60)
		time.Sleep(15 * time.Second)
		if readResp.LifecycleState == "SUCCEEDED" || readResp.LifecycleState == "FAILED" {
			deleteStack(ctx, stackID, client, solution)
			fmt.Printf("Delete complete successfully")
			break
		} else if readResp.LifecycleState == "FAILED" {
			fmt.Printf("Delete failed")
			break
		}
	}

	return *destroyJobResp.Job.Id

	/*
	   	readResp, err := client.GetJob(ctx, jobLifecycle)

	   	fmt.Println(readResp.LifecycleState)
	   	time.Sleep(15 * time.Second)
	   	fmt.Println(readResp2.LifecycleState)

	   	return *destroyJobResp.Job.Id
	   }
	*/
	/*
	   	getJobLifecycle := func() (interface{}, error) {
	   		request := resourcemanager.GetJobRequest{
	   			JobId: destroyJobResp.Id,
	   		}

	   		readResp, err := client.GetJob(ctx, request)

	   		if err != nil {
	   			return nil, err
	   		}

	   		return readResp.LifecycleState, err
	   	}

	   	fmt.Println(getJobLifecycle())
	   	time.Sleep(15 * time.Second)
	   	fmt.Println(getJobLifecycle())
	   	return *destroyJobResp.Job.Id
	   }

	   /*
	   	for {
	   		elapsed := int(time.Since(start).Seconds())
	   		destroyJobStatus := destroyJobResp.LifecycleState
	   		fmt.Printf("Current job status: %s [%dmin %dsec]\n", destroyJobStatus, elapsed/60, elapsed%60)
	   		time.Sleep(15 * time.Second)
	   		if destroyJobStatus == "SUCCEEDED" || destroyJobStatus == "FAILED" {
	   			fmt.Printf("Delete finished with status: %s", destroyJobStatus)
	   			break
	   		}
	   	}
	*/
}
