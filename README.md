# OCI Easy HPC deployment tool - ocihpc

`ocihpc` is a tool for simplifying deployments of HPC applications in Oracle Cloud Infrastructure (OCI).

## Prerequisites
The OCI user account you use in `ocihpc` should have the necessary policies configured for OCI Resource Manager. Please check [this link](https://docs.cloud.oracle.com/en-us/iaas/Content/Identity/Tasks/managingstacksandjobs.htm) for information on required policies.

## Installing ocihpc
### Installing ocihpc on Linux

1. Download the latest release with the following command:
```sh
curl -LO 
```

2. Make the ocihpc binary executable.
```sh
chmod +x ./ocihpc 
```

3. Move the ocihpc binary to your PATH.
```sh
sudo mv ./ocihpc /usr/local/bin/ocihpc 
```

4. Test that it works.
```sh
ocihpc version 
```

### Installing ocihpc on macOS

1. Download the latest release with the following command:
```sh
curl -LO 
```

2. Make the ocihpc binary executable.
```sh
chmod +x ./ocihpc 
```

3. Move the ocihpc binary to your PATH.
```sh
sudo mv ./ocihpc /usr/local/bin/ocihpc 
```

4. Test that it works.
```sh
ocihpc version 
```

### Installing ocihpc on Windows


1. Download the latest release with the following command:
```sh
curl -LO 
```

2. Add the ocihpc binary in to your PATH.

3. Test that it works.
```sh
ocihpc version 
```




## Using ocihpc

### 1 - List
You can get the list of available stacks by running `ocihpc list`.

Example:

```sh
$ ocihpc list

List of available stacks:

ClusterNetwork
Gromacs
OpenFOAM
```

### 2 - Initialize
Create a folder that you will use as the deployment source.

IMPORTANT: Use a different folder per stack. Do not initialize more than one stack in the same folder. Otherwise, the tool will overwrite the previous one.

Change to that folder and run `ocihpc init <stack name>`. `ocihpc` will download the necessary files to that folder.


```
$ mkdir ocihpc-test
$ cd ocihpc-test
$ ocihpc init ClusterNetwork

Downlading stack: ClusterNetwork

stack ClusterNetwork downloaded to /Users/opastirm/ocihpc-test/

IMPORTANT: Edit the contents of the /Users/opastirm/ocihpc-test/config.json file before running ocihpc deploy command
```

### 3 - Deploy
Before deploying, you need to change the values in `config.json` file. The variables depend on the stack you deploy. An example `config.json` for Cluster Network would look like this:

```json
{
  "variables": {
    "ad": "kWVD:PHX-AD-1",
    "bastion_ad": "kWVD:PHX-AD-2",
    "bastion_shape": "VM.Standard2.1",
    "node_count": "2",
    "ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAA......W6 opastirm@opastirm-mac"
  }
}
```

After you change the values in `config.json`, you can deploy the stack with `ocihpc deploy <stack name>`. This command will create a Stack on Oracle Cloud Resource Manager and deploy the stack using it.

For supported stacks, you can set the number of nodes you want to deploy by adding it to the `ocihpc deploy` command. If the stack does not support it or if you don't provide a value, the tool will deploy with the default numbers. 

For example, the following command will deploy a Cluster Network with 5 nodes:

```
$ ocihpc deploy ClusterNetwork 5
```

INFO: The tool will generate a deployment name that consists of `<stack name>-<current directory>-<random-number>`.

Example:

```
$ ocihpc deploy ClusterNetwork

Starting deployment...

Deploying ClusterNetwork-ocihpc-test-7355 [0min 0sec]
Deploying ClusterNetwork-ocihpc-test-7355 [0min 17sec]
Deploying ClusterNetwork-ocihpc-test-7355 [0min 35sec]
...
```

TIP: When running the `ocihpc deploy <stack name>` command, your shell might autocomplete it to the name of the zip file in the folder. This is fine. The tool will correct it, you don't need to delete the .zip extension from the command.

For example, `ocihpc deploy ClusterNetwork` and `ocihpc deploy ClusterNetwork.zip` are both valid commands.


### 4 - Connect
When deployment is completed, you will see the the bastion/headnode IP that you can connect to:

```
Successfully deployed ClusterNetwork-ocihpc-test-7355

You can connect to your head node using the command: ssh opc@$123.221.10.8 -i <location of the private key you used>

You can also find the IP address of the bastion/headnode in ClusterNetwork-ocihpc-test-7355_access.info file
```

### 5 - Delete
When you are done with your deployment, you can delete it by changing to the stack folder and running `ocihpc delete <stack name>`.

Example:
```
$ ocihpc delete ClusterNetwork

Deleting ClusterNetwork-ocihpc-test-7355 [0min 0sec]
Deleting ClusterNetwork-ocihpc-test-7355 [0min 17sec]
Deleting ClusterNetwork-ocihpc-test-7355 [0min 35sec]
...

Succesfully deleted ClusterNetwork-ocihpc-test-7355
```
