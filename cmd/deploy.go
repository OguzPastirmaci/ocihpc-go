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

		provider := common.DefaultConfigProvider()
		client, err := resourcemanager.NewResourceManagerClientWithConfigurationProvider(provider)
		helpers.FatalIfError(err)

		ctx := context.Background()

		stackID := createStack(ctx, provider, client, compartmentID, region, solution)
		//applyJobID := createApplyJob(ctx, provider, client, stackID, region)

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
}

func createStack(ctx context.Context, provider common.ConfigurationProvider, client resourcemanager.ResourceManagerClient, compartment string, region, solution string) string {
	stackName := fmt.Sprintf("%s-%s", solution, helpers.GetRandomString(4))
	//tenancyOcid, _ := provider.TenancyOCID()
	//compartmentID = os.Getenv("OCI_COMPARTMENT_ID")
	//region, _ := common.String(region)

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

	fmt.Println("Stack creation completed")
	return *stackResp.Stack.Id
}

func createApplyJob(ctx context.Context, provider common.ConfigurationProvider, client resourcemanager.ResourceManagerClient, stackID string, region string) string {

	jobReq := resourcemanager.CreateJobRequest{
		CreateJobDetails: resourcemanager.CreateJobDetails{
			StackId:   common.String(stackID),
			Operation: "APPLY",
			JobOperationDetails: resourcemanager.CreateApplyJobOperationDetails{
				ExecutionPlanStrategy: "AUTO_APPROVED",
			},
		},
	}

	jobResp, err := client.CreateJob(ctx, jobReq)

	if err != nil {
		fmt.Println("Submission of apply job failed", err)
		os.Exit(1)
	}

	fmt.Println("Apply job creation completed")
	return *jobResp.Job.Id
}

func writeStackInfo(key string, value string) {

	in := fmt.Sprintf("%s"+"="+"%s\n", key, value)

	f, err := os.OpenFile("stack.info", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(in)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
