// This software is licensed under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a stack for deployment",
	Long: `
Example command: ocihpc init --stack ClusterNetwork
	`,
	Run: func(cmd *cobra.Command, args []string) {
		stack, _ := cmd.Flags().GetString("stack")
		fmt.Printf("\nDownloading stack %s...", stack)
		stackInit(stack)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("stack", "s", "", "Name of the stack you want to deploy.")
	initCmd.MarkFlagRequired("stack")
}

func stackInit(stack string) {

	configURL := fmt.Sprintf("https://raw.githubusercontent.com/oracle-quickstart/oci-ocihpc/master/packages/%s/config.json", stack)
	zipURL := fmt.Sprintf("https://github.com/oracle-quickstart/oci-ocihpc/raw/master/packages/%s/%s.zip", stack, stack)

	configFilePath := getWd() + "/config.json"
	zipFilePath := getWd() + "/" + stack + ".zip"

	errConfig := downloadFile(configFilePath, configURL)
	if errConfig != nil {
		panic(errConfig)
	}

	errZip := downloadFile(zipFilePath, zipURL)
	if errZip != nil {
		panic(errZip)
	}
	fmt.Println("\n\nDownloaded stack " + stack)
	fmt.Printf("\nIMPORTANT: Edit the contents of the %s file before running ocihpc deploy command\n\n", configFilePath)
}
