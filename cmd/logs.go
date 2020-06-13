// This software is licensed under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/example/helpers"
	"github.com/oracle/oci-go-sdk/resourcemanager"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Get the logs of the last job",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		provider := common.DefaultConfigProvider()
		client, err := resourcemanager.NewResourceManagerClientWithConfigurationProvider(provider)
		helpers.FatalIfError(err)

		ctx := context.Background()
		jobID := getJobID()

		logs, _ := getTFLogs(ctx, provider, client, jobID)
		fmt.Println(logs)
	},
}

func init() {
	getCmd.AddCommand(logsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getTFLogs(ctx context.Context, provider common.ConfigurationProvider, client resourcemanager.ResourceManagerClient, jobID string) (string, error) {

	tf := resourcemanager.GetJobLogsRequest{
		JobId:                         &jobID,
		TimestampGreaterThanOrEqualTo: &common.SDKTime{time.Now().Add(time.Second * -300)},
		SortOrder:                     "ASC",
	}

	resp, err := client.GetJobLogs(ctx, tf)
	helpers.FatalIfError(err)

	logs, err := json.MarshalIndent(resp.Items, "", "    ")
	helpers.FatalIfError(err)

	return string(logs), err

}
