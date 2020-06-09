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

		provider := common.DefaultConfigProvider()
		client, err := resourcemanager.NewResourceManagerClientWithConfigurationProvider(provider)
		helpers.FatalIfError(err)

		ctx := context.Background()

		stackToDelete := getStackInfo("StackID")
		fmt.Println(stackToDelete)

		deleteStack(ctx, stackToDelete, client)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

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

	// delete a resource manager stack
	_, err := client.DeleteStack(ctx, req)
	helpers.FatalIfError(err)

	fmt.Println("Stack deletion")
}

func getStackInfo(value string) string {
	// Get substring after a string.

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
