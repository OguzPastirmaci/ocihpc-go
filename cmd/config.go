// This software is licensed under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl

package cmd

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/example/helpers"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure ocihpc",
	Long: `
Example command: ocihpc configure
	`,
	Run: func(cmd *cobra.Command, args []string) {

		//home, err := homedir.Dir()
		//helpers.FatalIfError(err)

		//configfile := home + "/.oci/config"

		provider := common.DefaultConfigProvider()

		if ok, _ := common.IsConfigurationProviderValid(provider); !ok {
			//fmt.Errorf("Did not find a valid configuration file. Answer the following questions to create one:", err)
			fmt.Printf("\nDid not find a valid configuration file. Answer the following questions to create one:\n\n")
			createKeys()
		} else {
			fmt.Printf("\nFound existing valid configuration. Exiting configuration.\n\n")
		}

		//createKeys()

	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

}

func createNewConfig() {

	//home, err := homedir.Dir()
	//helpers.FatalIfError(err)

	fmt.Println(`"The following links explain where to find the information required by this steps:
	User API Signing Key, OCID and Tenancy OCID:
	
    https://docs.cloud.oracle.com/Content/API/Concepts/apisigningkey.htm#Other
	
	Region:
    
    https://docs.cloud.oracle.com/Content/General/Concepts/regions.htm
	
	General config documentation:
    
    https://docs.cloud.oracle.com/Content/API/Concepts/sdkconfig.htm
    
"`)

	fmt.Printf("Enter a user OCID: \n")
	var userID string
	fmt.Scanln(&userID)

	fmt.Printf("Enter a tenancy OCID: \n")
	var tenancyID string
	fmt.Scanln(&tenancyID)

	fmt.Printf("Enter a region: \n")
	var region string
	fmt.Scanln(&region)

}

func createKeys() {

	home, err := homedir.Dir()
	helpers.FatalIfError(err)

	privateFileName := home + "/.oci/ocihpc_key.pem"
	publicFileName := home + "/.oci/ocihpc_key_public.pem"

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	helpers.FatalIfError(err)

	publicKey := key.PublicKey

	// Create private key
	outFile, err := os.Create(privateFileName)
	helpers.FatalIfError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	err = os.Chmod(privateFileName, 0600)
	helpers.FatalIfError(err)

	// Create public key
	bytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	helpers.FatalIfError(err)

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: bytes,
	}

	pemfile, err := os.Create(publicFileName)
	helpers.FatalIfError(err)
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	helpers.FatalIfError(err)

	md5sum := md5.Sum(pemkey.Bytes)

	hexarray := make([]string, len(md5sum))
	for i, c := range md5sum {
		hexarray[i] = hex.EncodeToString([]byte{c})
	}

	a := strings.Join(hexarray, ":")
	fmt.Println(a)
	//return strings.Join(hexarray, ":")

}

func colonSeparatedString(fingerprint [sha1.Size]byte) string {
	spaceSeparated := fmt.Sprintf("% x", fingerprint)
	return strings.Replace(spaceSeparated, " ", ":", -1)
}
