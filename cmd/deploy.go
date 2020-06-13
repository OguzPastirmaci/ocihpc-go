// This software is licensed under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl

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
	"time"

	"github.com/jeffail/gabs"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/example/helpers"
	"github.com/oracle/oci-go-sdk/resourcemanager"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:     "deploy",
	Aliases: []string{"create"},
	Short:   "Deploy a stack",
	Long: `
Example command: ocihpc deploy --stack ClusterNetwork --node-count 2 --region us-ashburn-1 --compartment-id ocid1.compartment.oc1..nus3q
	`,
	Run: func(cmd *cobra.Command, args []string) {
		stack, _ := cmd.Flags().GetString("stack")
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

		stackID := createStack(ctx, provider, client, compartmentID, region, stack, nodeCount)
		addStackInfo(stackID)

		applyJobID := createApplyJob(ctx, provider, client, stackID, region, stack)
		addJobInfo(applyJobID)

	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringP("compartment-id", "c", "", "Unique identifier (OCID) of the compartment that the stack will be deployed in.")
	deployCmd.MarkFlagRequired("compartment-id")

	deployCmd.Flags().StringP("region", "r", "", "The region to deploy to")

	deployCmd.Flags().StringP("stack", "s", "", "Name of the stack you want to deploy.")
	deployCmd.MarkFlagRequired("stack")

	deployCmd.Flags().StringP("node-count", "n", "", "Number of nodes to deploy.")
}

func createStack(ctx context.Context, provider common.ConfigurationProvider, client resourcemanager.ResourceManagerClient, compartment string, region string, stack string, nodeCount string) string {
	stackName := fmt.Sprintf("%s-%s", stack, helpers.GetRandomString(4))
	tenancyID, _ := provider.TenancyOCID()

	// Base64 the zip file
	zipFilePath := getWd() + "/" + stack + ".zip"
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

	config["tenancy_ocid"] = tenancyID
	config["compartment_ocid"] = compartment

	// region override
	_, r := config["region"]
	if r {
		if len(region) > 0 {
			config["region"] = region
		}
	} else {
		config["region"] = region
	}

	// node count override
	_, nc := config["node_count"]
	if nc {
		if len(nodeCount) > 0 {
			config["node_count"] = nodeCount
		}
	} else {
		fmt.Printf("\nChanging the node count is not supported with the stack %s, deploying stack with defaults.\n", stack)
	}

	req := resourcemanager.CreateStackRequest{
		CreateStackDetails: resourcemanager.CreateStackDetails{
			CompartmentId: common.String(compartment),
			ConfigSource: resourcemanager.CreateZipUploadConfigSourceDetails{
				ZipFileBase64Encoded: common.String(encoded),
			},
			DisplayName:      common.String(stackName),
			Description:      common.String(fmt.Sprintf("%s - Created by ocihpc", stack)),
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

func createApplyJob(ctx context.Context, provider common.ConfigurationProvider, client resourcemanager.ResourceManagerClient, stackID string, region string, stack string) string {

	outputQuery := map[string]string{
		"ClusterNetwork": "outputs.bastion.value",
		"vcn":            "outputs.vcn_id.value",
	}

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
		fmt.Println("Deployment failed with the following errors:\n\n", err)
		os.Exit(1)
	}

	jobLifecycle := resourcemanager.GetJobRequest{
		JobId: applyJobResp.Id,
	}

	fmt.Println()
	start := time.Now().Add(time.Second * -5)

	for {
		elapsed := int(time.Since(start).Seconds())
		readResp, err := client.GetJob(ctx, jobLifecycle)

		if err != nil {
			fmt.Println("Deployment failed with the following errors:\n\n", err)
			os.Exit(1)
		}

		fmt.Printf("Deploying stack: %s [%dmin %dsec]\n", stack, elapsed/60, elapsed%60)
		time.Sleep(10 * time.Second)

		if readResp.LifecycleState == "SUCCEEDED" {
			fmt.Printf("\nDeployment of %s completed successfully\n", stack)

			tfStateReq := resourcemanager.GetJobTfStateRequest{
				JobId: applyJobResp.Id,
			}
			tfStateResp, _ := client.GetJobTfState(ctx, tfStateReq)
			body, _ := ioutil.ReadAll(tfStateResp.Content)
			tfStateParsed, err := gabs.ParseJSON([]byte(string(body)))
			if err != nil {
				log.Fatal("Error:", err)
			}
			var outputIP string
			outputIP = tfStateParsed.Path(outputQuery[stack]).Data().(string)
			fmt.Printf("\nYou can connect to your bastion/headnode using the command: ssh opc@%s -i <location of the private key>\n\n", outputIP)
			break
		} else if readResp.LifecycleState == "FAILED" {
			fmt.Printf("\nDeployment failed. You can run 'ocihpc get logs' to get the logs of the failed job\n")
			fmt.Printf("\nPlease note that there might be some resources that are already created. Run 'ocihpc delete %s' to delete those resources.", stack)
			break
		}
	}

	return *applyJobResp.Job.Id
}
