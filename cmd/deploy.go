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
		createStackInfo(".solution.json", stackID, region)

		createApplyJob(ctx, provider, client, stackID, region, solution)

	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringP("compartment-id", "c", "", "Unique identifier (OCID) of the compartment in which the stack resides.")
	deployCmd.MarkFlagRequired("compartment-id")

	deployCmd.Flags().StringP("region", "r", "", "The region to deploy to")

	deployCmd.Flags().StringP("solution", "s", "", "Name of the solution you want to deploy.")
	deployCmd.MarkFlagRequired("solution")

	deployCmd.Flags().StringP("node-count", "n", "", "Number of nodes to deploy.")
}

func createStack(ctx context.Context, provider common.ConfigurationProvider, client resourcemanager.ResourceManagerClient, compartment string, region string, solution string, nodeCount string) string {
	stackName := fmt.Sprintf("%s-%s", solution, helpers.GetRandomString(4))
	//tenancyOcid, _ := provider.TenancyOCID()
	//compartmentID = os.Getenv("OCI_COMPARTMENT_ID")

	// Base64 the zip file
	zipFilePath := getWd() + "/" + solution + ".zip"
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

	// node count override
	_, nc := config["node_count"]
	if nc {
		if len(nodeCount) > 0 {
			config["node_count"] = nodeCount
			fmt.Println("Changing node count.")
		}
	} else {
		fmt.Printf("\nChanging the node count is not supported with the solution %s, deploying with defaults.\n", solution)
	}

	// region override
	_, r := config["region"]
	if r {
		if len(region) > 0 {
			config["region"] = region
			fmt.Println("Changing region.")
		}
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

func createApplyJob(ctx context.Context, provider common.ConfigurationProvider, client resourcemanager.ResourceManagerClient, stackID string, region string, solution string) string {

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

	start := time.Now()

	for {
		elapsed := int(time.Since(start).Seconds())
		readResp, err := client.GetJob(ctx, jobLifecycle)

		if err != nil {
			fmt.Println("Deployment failed with the following errors:\n\n", err)
			os.Exit(1)
		}

		fmt.Printf("Deploying solution: %s [%dmin %dsec]\n", solution, elapsed/60, elapsed%60)
		time.Sleep(10 * time.Second)

		if readResp.LifecycleState == "SUCCEEDED" {
			fmt.Printf("Deployment completed successfully\n")

			tfStateReq := resourcemanager.GetJobTfStateRequest{
				JobId: applyJobResp.Id,
			}
			tfStateResp, _ := client.GetJobTfState(ctx, tfStateReq)
			body, _ := ioutil.ReadAll(tfStateResp.Content)
			tfStateParsed, err := gabs.ParseJSON([]byte(string(body)))
			if err != nil {
				log.Fatal("Error:", err)
			}
			var bastionIP string
			bastionIP = tfStateParsed.Path("outputs.bastion.value").Data().(string)
			fmt.Printf("\nYou can connect to your head node using the command: ssh opc@%s -i <location of the private key you used>\n", bastionIP)
			break
		} else if readResp.LifecycleState == "FAILED" {
			fmt.Printf("\nDeployment failed. Please note there might be some resources already created.\n")
			fmt.Printf("\nRun ocihpc delete %s to delete those resources.", solution)
			break
		}
	}

	return *applyJobResp.Job.Id
}
