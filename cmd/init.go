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
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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

func pwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func downloadFile(filepath string, url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func solutionInit(solution_name string) {

	configUrl := fmt.Sprintf("https://raw.githubusercontent.com/oracle-quickstart/oci-ocihpc/master/packages/%s/config.json", solution_name)
	zipUrl := fmt.Sprintf("https://github.com/oracle-quickstart/oci-ocihpc/raw/master/packages/%s/%s.zip", solution_name, solution_name)

	configFilePath := pwd() + "/config.json"
	zipFilePath := pwd() + "/" + solution_name + ".zip"

	errConfig := downloadFile(configFilePath, configUrl)
	if errConfig != nil {
		panic(errConfig)
	}

	errZip := downloadFile(zipFilePath, zipUrl)
	if errZip != nil {
		panic(errZip)
	}
	fmt.Println("\n\nDownloaded solution " + solution_name)
	fmt.Printf("\nIMPORTANT: Edit the contents of the %s file before running ocihpc deploy command\n\n", configFilePath)
}
