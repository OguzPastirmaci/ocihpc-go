// This software is licensed under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Get the IP address of the headnode or the bastion of a stack",
	Long: `
Example command: ocihpc get ip
	`,

	Run: func(cmd *cobra.Command, args []string) {

		stackIP := getStackIP()
		stackName := getStackName()
		fmt.Printf("\nYou can connect to your bastion/headnode using the following command:\n\n")
		fmt.Printf("ssh %s@%s -i <location of the private key>\n\n", stackUser[stackName], stackIP)

	},
}

func init() {
	getCmd.AddCommand(ipCmd)
}
