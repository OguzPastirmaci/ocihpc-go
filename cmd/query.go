// This software is licensed under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl

package cmd

var outputQuery = map[string]string{
	"ClusterNetwork": "outputs.bastion.value",
	"vcn":            "outputs.vcn_id.value",
}

var stackUser = map[string]string{
	"ClusterNetwork": "opc",
}

var stackVersion = map[string]string{
	"ClusterNetwork": "0.12.x",
	"vcn":            "0.12.x",
}
