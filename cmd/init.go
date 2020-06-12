// This software is licensed under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		solution, _ := cmd.Flags().GetString("solution")
		fmt.Printf("\nDownloading solution %s...", solution)
		solutionInit(solution)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("solution", "s", "", "Name of the solution you want to deploy.")
	initCmd.MarkFlagRequired("solution")
}

func solutionInit(solution string) {

	configURL := fmt.Sprintf("https://raw.githubusercontent.com/oracle-quickstart/oci-ocihpc/master/packages/%s/config.json", solution)
	zipURL := fmt.Sprintf("https://github.com/oracle-quickstart/oci-ocihpc/raw/master/packages/%s/%s.zip", solution, solution)

	configFilePath := getWd() + "/config.json"
	zipFilePath := getWd() + "/" + solution + ".zip"

	errConfig := downloadFile(configFilePath, configURL)
	if errConfig != nil {
		panic(errConfig)
	}

	errZip := downloadFile(zipFilePath, zipURL)
	if errZip != nil {
		panic(errZip)
	}
	fmt.Println("\n\nDownloaded solution " + solution)
	fmt.Printf("\nIMPORTANT: Edit the contents of the %s file before running ocihpc deploy command\n\n", configFilePath)
}
