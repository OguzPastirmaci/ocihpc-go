// This software is licensed under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl

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

		provider := common.DefaultConfigProvider()
		solution, _ := cmd.Flags().GetString("solution")
		client, err := resourcemanager.NewResourceManagerClientWithConfigurationProvider(provider)
		if err != nil {
			panic(err)
		}
		helpers.FatalIfError(err)

		ctx := context.Background()
		stackID := getJSON(".solution.json", "stack_info.stackID")

		createDestroyJob(ctx, provider, client, stackID, solution)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringP("solution", "s", "", "Solution to delete")
	deleteCmd.MarkFlagRequired("solution")
}

func deleteStack(ctx context.Context, stackID string, client resourcemanager.ResourceManagerClient, solution string) {

	req := resourcemanager.DeleteStackRequest{
		StackId: common.String(stackID),
	}

	_, err := client.DeleteStack(ctx, req)
	helpers.FatalIfError(err)
}

func createDestroyJob(ctx context.Context, provider common.ConfigurationProvider, client resourcemanager.ResourceManagerClient, stackID string, solution string) string {

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

	if err != nil {
		fmt.Println("Delete failed with the following errors:\n\n", err)
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
			fmt.Println("Delete failed with the following errors:\n\n", err)
			os.Exit(1)
		}

		fmt.Printf("Deleting solution: %s [%dmin %dsec]\n", solution, elapsed/60, elapsed%60)
		time.Sleep(15 * time.Second)
		if readResp.LifecycleState == "SUCCEEDED" {
			deleteStack(ctx, stackID, client, solution)
			fmt.Printf("Delete completed successfully")
			os.Remove("stack.info")
			break
		} else if readResp.LifecycleState == "FAILED" {
			fmt.Printf("Delete failed")
			break
		}
	}

	return *destroyJobResp.Job.Id
}
