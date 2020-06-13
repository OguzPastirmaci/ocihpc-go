// This software is licensed under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl

package cmd

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/oracle/oci-go-sdk/example/helpers"
)

var filename = ".stackinfo.json"

func addStackInfo(s Stack) {

	file, _ := json.MarshalIndent(s, "", " ")

	_ = ioutil.WriteFile(filename, file, 0644)
}

func getStackName() string {

	content, err := ioutil.ReadFile(filename)
	helpers.FatalIfError(err)

	var info Stack
	json.Unmarshal([]byte(content), &info)

	return info.StackName
}

func getStackID() string {

	content, err := ioutil.ReadFile(filename)
	helpers.FatalIfError(err)

	var info Stack
	json.Unmarshal([]byte(content), &info)

	return info.StackID
}

func getStackIP() string {

	content, err := ioutil.ReadFile(filename)
	helpers.FatalIfError(err)

	var info Stack
	json.Unmarshal([]byte(content), &info)

	return info.StackIP
}

func getJobID() string {

	content, err := ioutil.ReadFile(filename)
	helpers.FatalIfError(err)

	var info Stack
	json.Unmarshal([]byte(content), &info)

	return info.JobID
}

func getWd() string {
	dir, err := os.Getwd()
	helpers.FatalIfError(err)

	return dir
}

func downloadFile(filepath string, url string) error {

	resp, err := http.Get(url)
	helpers.FatalIfError(err)

	defer resp.Body.Close()

	out, err := os.Create(filepath)
	helpers.FatalIfError(err)

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
