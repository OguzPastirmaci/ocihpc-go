// This software is licensed under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		stackIP := getStackIP()
		stackName := getStackName()
		fmt.Println("You can connect to your bastion/headnode using the following command:")
		fmt.Printf("ssh %s@%s -i <location of the private key>\n\n", stackUser[stackName], stackIP)

	},
}

func init() {
	getCmd.AddCommand(ipCmd)
}
