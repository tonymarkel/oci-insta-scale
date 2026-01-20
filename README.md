![Insta-Scale](images/insta-scale.png)
# OCI Instance Pool Manager
_Generated using an LLM_

A Go program to create, scale, and manage Oracle Cloud Infrastructure (OCI) instance pools with an arbitrary number of virtual machines.

## Features

- ✅ Create instance pools with configurable VM count
- ✅ Scale existing instance pools up or down
- ✅ Terminate instance pools
- ✅ Support for flexible shapes with custom OCPU and memory
- ✅ Load balancer integration
- ✅ Multi-availability domain placement
- ✅ Fault domain distribution
- ✅ Custom SSH keys, user data, and metadata
- ✅ Tagging support (freeform and defined tags)

## Prerequisites

1. **OCI Account**: An active Oracle Cloud Infrastructure account
2. **Go**: Go 1.21 or later installed
3. **OCI CLI Configuration**: Your OCI credentials configured

### OCI Authentication Setup

You need the following information from your OCI account:

- **Tenancy OCID**: Your OCI tenancy ID
- **User OCID**: Your user ID
- **API Key Fingerprint**: From your API key
- **Private Key Path**: Path to your API private key file
- **Region**: OCI region (e.g., `us-phoenix-1`)

You can find these in the OCI Console under your user profile settings.

## Installation

1. Clone this repository:
```bash
git clone https://github.com/tonymarkel/oci-insta-scale
cd oci-insta-scale
```

2. Install dependencies:
```bash
go mod download
```

3. Build the program:
```bash
go build -o oci-insta-scale
```

## Configuration

Create a `config.yaml` file with your OCI settings. See `config.example.yaml` for a complete example.

### Minimal Configuration

```yaml
# OCI Authentication
tenancy_ocid: "ocid1.tenancy.oc1..aaaaa..."
user_ocid: "ocid1.user.oc1..aaaaa..."
fingerprint: "aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99"
private_key_path: "/path/to/.oci/oci_api_key.pem"
region: "us-phoenix-1"

# Compartment
compartment_id: "ocid1.compartment.oc1..aaaaa..."

# Instance Pool Configuration
instance_pool:
  display_name: "my-instance-pool"
  size: 3  # Number of VMs
  
  instance_configuration:
    display_name: "my-instance-config"
    shape: "VM.Standard.E4.Flex"
    shape_config:
      ocpus: 1
      memory_in_gbs: 16
    image_id: "ocid1.image.oc1.phx.aaaaa..."  # OS image
    subnet_id: "ocid1.subnet.oc1.phx.aaaaa..."
    assign_public_ip: true
    ssh_authorized_keys: "ssh-rsa AAAAB3NzaC1yc2E..."
  
  placement:
    - availability_domain: "IYiP:PHX-AD-1"
```

## Usage

### Create an Instance Pool

Create a new instance pool with the number of VMs specified in the config:

```bash
./oci-insta-scale -config config.yaml -action create
```

Override the instance count from command line:

```bash
./oci-insta-scale -config config.yaml -action create
```

### Scale an Existing Instance Pool

Scale an instance pool to a different size:

```bash
./oci-insta-scale -config config.yaml -action scale \
  -pool-id ocid1.instancepool.oc1.phx.aaaaa... \
```

### Terminate an Instance Pool

Terminate an instance pool and all its instances:

```bash
./oci-insta-scale -config config.yaml -action terminate \
  -pool-id ocid1.instancepool.oc1.phx.aaaaa...
```

## Command-Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-config` | Path to configuration file | `config.yaml` |
| `-action` | Action to perform: `create`, `scale`, `terminate` | `create` |
| `-count` | Number of instances (overrides config) | 0 |
| `-compartment` | Compartment OCID (overrides config) | "" |
| `-name` | Instance pool display name (overrides config) | "" |
| `-pool-id` | Instance pool ID (for scale/terminate) | "" |

## Advanced Configuration

### Flexible Shapes

For flexible shapes, specify OCPUs and memory:

```yaml
instance_configuration:
  shape: "VM.Standard.E4.Flex"
  shape_config:
    ocpus: 2
    memory_in_gbs: 32
```

### Multiple Availability Domains

Distribute instances across multiple ADs:

```yaml
placement:
  - availability_domain: "IYiP:PHX-AD-1"
    fault_domains: ["FAULT-DOMAIN-1", "FAULT-DOMAIN-2"]
  - availability_domain: "IYiP:PHX-AD-2"
  - availability_domain: "IYiP:PHX-AD-3"
```

### Load Balancer Integration

Attach instances to a load balancer:

```yaml
instance_pool:
  load_balancers:
    - load_balancer_id: "ocid1.loadbalancer.oc1.phx.aaaaa..."
      backend_set_name: "backend-set-1"
      port: 80
      vnic_selection: "PrimaryVnic"
```

### User Data and Cloud-Init

Provide cloud-init script (base64 encoded):

```yaml
instance_configuration:
  user_data: "IyEvYmluL2Jhc2gKZWNobyAiSGVsbG8gV29ybGQi"
```

### Custom Metadata

Add custom metadata to instances:

```yaml
instance_configuration:
  metadata:
    application: "web-server"
    environment: "production"
```

### Tagging

Add freeform or defined tags:

```yaml
instance_configuration:
  freeform_tags:
    Project: "MyProject"
    CostCenter: "12345"
  defined_tags:
    Operations:
      Environment: "Production"
```

## Examples

### Example 1: Simple Web Server Pool

Create a pool of 5 web servers:

```bash
./oci-insta-scale -config web-config.yaml -action create -count 5
```

### Example 2: Scale Up During Peak Hours

Scale from 3 to 10 instances:

```bash
./oci-insta-scale -config config.yaml -action scale \
  -pool-id ocid1.instancepool.oc1.phx.aaaaa... \
  -count 10
```

### Example 3: Scale Down After Peak

Scale back down to 3 instances:

```bash
./oci-insta-scale -config config.yaml -action scale \
  -pool-id ocid1.instancepool.oc1.phx.aaaaa... \
  -count 3
```

## Finding Required OCIDs

### Image OCID

List available images in your region:

```bash
oci compute image list --compartment-id <compartment-ocid> \
  --operating-system "Oracle Linux" --shape "VM.Standard.E4.Flex"
```

### Subnet OCID

List subnets in your VCN:

```bash
oci network subnet list --compartment-id <compartment-ocid> \
  --vcn-id <vcn-ocid>
```

### Availability Domains

List availability domains:

```bash
oci iam availability-domain list --compartment-id <compartment-ocid>
```

## Troubleshooting

### Authentication Errors

- Verify your API key is valid and uploaded to OCI
- Check that the private key file path is correct
- Ensure the fingerprint matches your API key

### Permission Errors

- Verify your user has permissions to create instance pools
- Check compartment access policies
- Ensure you have quota available for the requested resources

### Configuration Errors

The program validates your configuration and will report specific errors if required fields are missing.

## Project Structure

```
.
├── main.go           # Main entry point and CLI handling
├── config.go         # Configuration loading and validation
├── oci_client.go     # OCI SDK client wrapper and operations
├── go.mod            # Go module dependencies
├── config.yaml       # Your configuration file (create this)
└── README.md         # This file
```

## Dependencies

- `github.com/oracle/oci-go-sdk/v65` - Oracle Cloud Infrastructure Go SDK
- `gopkg.in/yaml.v3` - YAML configuration parsing

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues related to:
- **OCI SDK**: See [OCI Go SDK Documentation](https://docs.oracle.com/en-us/iaas/tools/go/latest/)
- **OCI Instance Pools**: See [OCI Instance Pools Documentation](https://docs.oracle.com/en-us/iaas/Content/Compute/Tasks/creatinginstancepool.htm)
