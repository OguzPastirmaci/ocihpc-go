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
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/example/helpers"
	"github.com/oracle/oci-go-sdk/resourcemanager"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
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
		client, err := resourcemanager.NewResourceManagerClientWithConfigurationProvider(provider)
		helpers.FatalIfError(err)

		ctx := context.Background()

		//stackToDelete := getStackInfo("StackID")
		//fmt.Println(stackToDelete)
		createDestroyJob(ctx, provider, client, "ocid1.ormstack.oc1.iad.aaaaaaaarrxefa7ogac7gu5fstm5fniblunh6gvzgtpezjkhheig3afdw67q", region)

		//deleteStack(ctx, stackToDelete, client)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringP("region", "r", "", "The region to deploy to")
	deleteCmd.MarkFlagRequired("region")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func deleteStack(ctx context.Context, stackID string, client resourcemanager.ResourceManagerClient) {

	req := resourcemanager.DeleteStackRequest{
		StackId: common.String(stackID),
	}

	_, err := client.DeleteStack(ctx, req)
	helpers.FatalIfError(err)

	fmt.Println("Stack deletion")
}

func createDestroyJob(ctx context.Context, provider common.ConfigurationProvider, client resourcemanager.ResourceManagerClient, stackID string, region string) string {

	jobReq := resourcemanager.CreateJobRequest{
		CreateJobDetails: resourcemanager.CreateJobDetails{
			StackId:   common.String(stackID),
			Operation: "DESTROY",
			JobOperationDetails: resourcemanager.CreateDestroyJobOperationDetails{
				ExecutionPlanStrategy: "AUTO_APPROVED",
			},
		},
	}

	jobResp, err := client.CreateJob(ctx, jobReq)

	if err != nil {
		fmt.Println("Submission of destroy job failed", err)
		os.Exit(1)
	}

	fmt.Println("Destroy job creation completed")
	return *jobResp.Job.Id
}

func getStackInfo(value string) string {

	a := value + "="

	content, err := ioutil.ReadFile("stack.info")
	if err != nil {
		log.Fatal(err)
	}

	text := string(content)

	pos := strings.LastIndex(text, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(text) {
		return ""
	}
	return text[adjustedPos:len(text)]
}