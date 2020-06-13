// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.

// Example code for sending raw request to  Service API

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/example/helpers"
)

// ExampleRawRequest compose a request, sign it and send to server

func main() {
	getTFState()

}

func getTFState() {
	// build the url

	jobId := "ocid1.ormjob.oc1.iad.aaaaaaaakyzp3s46mec66fssuzfp5u7nzorgil2thyajpzrliwb7xxqt73dq"
	url := "https://resourcemanager.us-ashburn-1.oraclecloud.com/20180917/jobs/" + jobId + "/logs"

	// create request
	request, err := http.NewRequest("GET", url, nil)
	helpers.FatalIfError(err)

	// Set the Date header
	request.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))

	// And a provider of cryptographic keys
	provider := common.DefaultConfigProvider()

	// Build the signer
	signer := common.DefaultRequestSigner(provider)

	// Sign the request
	signer.Sign(request)

	client := http.Client{}

	//fmt.Println("send request")

	// Execute the request
	resp, err := client.Do(request)
	helpers.FatalIfError(err)

	defer resp.Body.Close()

	//log.Println("response Status:", resp.Status)
	//log.Println("response Headers:", resp.Header)

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))

	//fmt.Println("receive response")
}

/*
func getTFLogs(ctx context.Context, provider common.ConfigurationProvider, client resourcemanager.ResourceManagerClient, jobID string) (string, error) {

	tf := resourcemanager.GetJobLogsRequest{
		JobId:                         &jobID,
		TimestampGreaterThanOrEqualTo: &common.SDKTime{time.Now().Add(time.Second * -300)},
	}

	resp, err := client.GetJobLogs(ctx, tf)
	helpers.FatalIfError(err)

	fmt.Println(resp)

	return resp.Items, err
}
*/
