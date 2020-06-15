/*
 * Genarate rsa keys.
 */

 package main

 import (
	 "crypto/rand"
	 "crypto/rsa"
	 "crypto/x509"
	 "encoding/pem"
	 "fmt"
	 "os"
	 "bufio"
 )
 
 func main() {

	/*
	reader := bufio.NewReader(os.Stdin)

fmt.Printf(`This command provides a walkthrough of creating a valid CLI config file.

The following links explain where to find the information required by this script:

User API Signing Key, OCID and Tenancy OCID:

	https://docs.cloud.oracle.com/Content/API/Concepts/apisigningkey.htm#Other

Region:

	https://docs.cloud.oracle.com/Content/General/Concepts/regions.htm

General config documentation:

	https://docs.cloud.oracle.com/Content/API/Concepts/sdkconfig.htm`)	



    fmt.Print("Enter your use OCID: ")
	user, _ := reader.ReadString('\n')
	
	fmt.Print("Enter your tenancy OCID: ")
	tenancy, _ := reader.ReadString('\n')


	fmt.Println(user, tenancy)

*/

	privateFileName := "ocihpc_key.pem"
	publicFileName := "ocihpc_key_public.pem"

	 key, err := rsa.GenerateKey(rand.Reader, 2048)
	 checkError(err)
 
	 publicKey := key.PublicKey
 
	 savePEMKey(privateFileName, key)
 
	 savePublicPEMKey(publicFileName, publicKey)
 }
 

 func createPEMKey(fileName string, key *rsa.PrivateKey) {
	 outFile, err := os.Create(fileName)
	 checkError(err)
	 defer outFile.Close()
 
	 var privateKey = &pem.Block{
		 Type:  "PRIVATE KEY",
		 Bytes: x509.MarshalPKCS1PrivateKey(key),
	 }
 
	 err = pem.Encode(outFile, privateKey)
	 err = os.Chmod(fileName, 0600)
	 checkError(err)
	
 }
 
 func createPublicPEMKey(fileName string, pubkey rsa.PublicKey) {
	bytes, err := x509.MarshalPKIXPublicKey(&pubkey)
	 checkError(err)
 
	 var pemkey = &pem.Block{
		 Type:  "PUBLIC KEY",
		 Bytes: bytes,
	 }
 
	 pemfile, err := os.Create(fileName)
	 checkError(err)
	 defer pemfile.Close()
 
	 err = pem.Encode(pemfile, pemkey)
	 checkError(err)
 }
 
 func checkError(err error) {
	 if err != nil {
		 fmt.Println("Fatal error ", err.Error())
		 os.Exit(1)
	 }
 }