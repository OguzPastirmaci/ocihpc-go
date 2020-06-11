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
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/example/helpers"
	"github.com/oracle/oci-go-sdk/resourcemanager"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "A brief description of your command",
	Long: `
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		solution, _ := cmd.Flags().GetString("solution")
		region, _ := cmd.Flags().GetString("region")
		compartmentID, _ := cmd.Flags().GetString("compartment-id")
		nodeCount, _ := cmd.Flags().GetString("node-count")

		if _, err := strconv.Atoi(nodeCount); err != nil {
			fmt.Printf("\nNode count must be a number, you entered: %s\n", nodeCount)
			os.Exit(1)
		}

		provider := common.DefaultConfigProvider()
		client, err := resourcemanager.NewResourceManagerClientWithConfigurationProvider(provider)
		if err != nil {
			panic(err)
		}

		helpers.FatalIfError(err)

		ctx := context.Background()

		stackID := createStack(ctx, provider, client, compartmentID, region, solution, nodeCount)
		//createApplyJob(ctx, provider, client, stackID, region)

		writeStackInfo("StackID", stackID)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringP("compartment-id", "c", "", "Unique identifier (OCID) of the compartment in which the stack resides.")
	deployCmd.MarkFlagRequired("compartment-id")

	deployCmd.Flags().StringP("region", "r", "", "The region to deploy to")
	deployCmd.MarkFlagRequired("region")

	deployCmd.Flags().StringP("solution", "s", "", "Name of the solution you want to deploy.")
	deployCmd.MarkFlagRequired("solution")

	deployCmd.Flags().StringP("node-count", "n", "", "Number of nodes to deploy.")
}

func createStack(ctx context.Context, provider common.ConfigurationProvider, client resourcemanager.ResourceManagerClient, compartment string, region, solution string, nodeCount string) string {
	stackName := fmt.Sprintf("%s-%s", solution, helpers.GetRandomString(4))
	//tenancyOcid, _ := provider.TenancyOCID()
	//compartmentID = os.Getenv("OCI_COMPARTMENT_ID")
	//region, _ := common.String(region)

	// Base64 the zip file
	zipFilePath := pwd() + "/" + solution + ".zip"
	f, _ := os.Open(zipFilePath)
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	encoded := base64.StdEncoding.EncodeToString(content)

	// read config.json
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var config map[string]string
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatal(err)
	}

	// if node count was entered, update the map

	_, ok := config["node_count"]
	if ok {
		if len(nodeCount) > 0 {
			config["node_count"] = nodeCount
			fmt.Println("Changing node count.")
		}
	} else {
		fmt.Printf("Changing the node count is not supported with the solution %s, deploying with defaults.", solution)
	}

	req := resourcemanager.CreateStackRequest{
		CreateStackDetails: resourcemanager.CreateStackDetails{
			CompartmentId: common.String(compartment),
			ConfigSource: resourcemanager.CreateZipUploadConfigSourceDetails{
				ZipFileBase64Encoded: common.String(encoded),
			},
			DisplayName:      common.String(stackName),
			Description:      common.String(fmt.Sprintf("%s - Created by ocihpc", solution)),
			Variables:        config,
			TerraformVersion: common.String("0.12.x"),
		},
	}

	stackResp, err := client.CreateStack(ctx, req)
	helpers.FatalIfError(err)

	if err != nil {
		fmt.Println("Stack creation failed: ", err)
		os.Exit(1)
	}

	return *stackResp.Stack.Id

}

func createApplyJob(ctx context.Context, provider common.ConfigurationProvider, client resourcemanager.ResourceManagerClient, stackID string, region string) string {

	applyJobReq := resourcemanager.CreateJobRequest{
		CreateJobDetails: resourcemanager.CreateJobDetails{
			StackId:   common.String(stackID),
			Operation: "APPLY",
			JobOperationDetails: resourcemanager.CreateApplyJobOperationDetails{
				ExecutionPlanStrategy: "AUTO_APPROVED",
			},
		},
	}

	applyJobResp, err := client.CreateJob(ctx, applyJobReq)

	if err != nil {
		fmt.Println("Submission of apply job failed", err)
		os.Exit(1)
	}

	//fmt.Println("Apply job creation completed")

	/*
		fmt.Printf("Stack ID of the apply job is: %s\n", stackID)
		fmt.Printf("Job ID of the apply job is: %s\n", *applyJobResp.Id)
		fmt.Printf("Lifecycle state of the apply job is: %s\n", applyJobResp.LifecycleState)
		fmt.Println("Waiting for 30 seconds before checking the state again")
		time.Sleep(30 * time.Second)
		fmt.Printf("Lifecycle state of the apply job is: %s\n", applyJobResp.LifecycleState)
	*/

	/*
		for {
			applyJobStatus := jobResp.LifecycleState
			fmt.Printf("Current job status: %s\n", applyJobStatus)
			//fmt.Printf("Current job status: %s\n", *jobResp.LifecycleState)
			time.Sleep(15 * time.Second)
			if applyJobStatus == "SUCCEEDED" || applyJobStatus == "FAILED" {
				fmt.Printf("Apply finished with status: %s", applyJobStatus)
				break
			}
		}
	*/
	return *applyJobResp.Job.Id
}
