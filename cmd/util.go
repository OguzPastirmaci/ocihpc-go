// This software is licensed under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl

package cmd

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/jeffail/gabs"
)

func getJSON(filename string, value string) string {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	jsonParsed, err := gabs.ParseJSON(data)
	result := jsonParsed.Path(value).Data().(string)
	return result
}

func getWd() string {
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

func addStackInfo(stackID string) {
	stackInfo := gabs.New()
	stackInfo.Set(stackID, "stack_info", "stackID")
	ioutil.WriteFile(".stack_info.json", []byte(stackInfo.StringIndent("", "  ")), 0644)
}

func addJobInfo(jobID string) {
	jobInfo := gabs.New()
	jobInfo.Set(jobID, "job_info", "jobID")
	ioutil.WriteFile(".job_info.json", []byte(jobInfo.StringIndent("", "  ")), 0644)
}
