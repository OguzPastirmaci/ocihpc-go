// This software is licensed under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl

package cmd

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
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

		home, err := homedir.Dir()
		helpers.FatalIfError(err)

		configfile := home + "/.oci/config"

		provider := common.DefaultConfigProvider()

		if ok, _ := common.IsConfigurationProviderValid(provider); !ok {
			//fmt.Errorf("Did not find a valid configuration file. Answer the following questions to create one:", err)
			fmt.Printf("\nDid not find a valid configuration file. Answer the following questions to create one:\n\n")
			createNewConfig(configfile)
		} else {
			fmt.Printf("\nFound existing valid configuration. Exiting configuration.\n\n")
		}

	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

}

func createNewConfig(configfile string) {

	home, err := homedir.Dir()
	helpers.FatalIfError(err)

	var user string
	var tenancy string
	var region string
	var fingerprint string

	privateFileName := home + "/.oci/ocihpc_key.pem"
	publicFileName := home + "/.oci/ocihpc_key_public.pem"

	file, err := os.OpenFile(configfile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
	helpers.FatalIfError(err)

	defer file.Close()

	fmt.Printf(`The following links explain where to find the information required by this steps:

	User API Signing Key, OCID and Tenancy OCID:
	
    https://docs.cloud.oracle.com/Content/API/Concepts/apisigningkey.htm#Other
	
	Region:
    
    https://docs.cloud.oracle.com/Content/General/Concepts/regions.htm
	
	General config documentation:
    
    https://docs.cloud.oracle.com/Content/API/Concepts/sdkconfig.htm
    
`)

	fmt.Printf("\nEnter a user OCID: ")
	fmt.Scanln(&user)

	fmt.Printf("\nEnter a tenancy OCID: ")
	fmt.Scanln(&tenancy)

	fmt.Printf("\nEnter a region: ")
	fmt.Scanln(&region)

	fingerprint = createKeys(privateFileName, publicFileName)

	content := fmt.Sprintf("[DEFAULT]\nuser=%s\nfingerprint=%s\nkey_file=%s\ntenancy=%s\nregion=%s", user, fingerprint, privateFileName, tenancy, region)
	_, err = file.WriteString(content)
}

func createKeys(privateFileName string, publicFileName string) string {

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	helpers.FatalIfError(err)

	publicKey := key.PublicKey

	// Create private key
	outFile, err := os.Create(privateFileName)
	helpers.FatalIfError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	err = os.Chmod(privateFileName, 0600)
	helpers.FatalIfError(err)

	// Create public key
	bytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	helpers.FatalIfError(err)

	var pemkey = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: bytes,
	}

	pemfile, err := os.Create(publicFileName)
	helpers.FatalIfError(err)
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	helpers.FatalIfError(err)

	md5sum := md5.Sum(pemkey.Bytes)

	fp := make([]string, len(md5sum))
	for i, c := range md5sum {
		fp[i] = hex.EncodeToString([]byte{c})
	}

	fingerprint := strings.Join(fp, ":")
	return fingerprint

}
