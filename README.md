# OCI Insta-Scale

A Go application to create an arbitrary number of Oracle Cloud Infrastructure (OCI) compute instances using capacity reservations.

## Features

- Create multiple OCI compute instances with a single command
- Support for capacity reservations
- Flexible shape configuration (OCPUs and memory)
- Configurable via YAML
- Dry-run mode for testing
- Progress tracking and status monitoring
- Automatic waiting for instances to reach RUNNING state

## Prerequisites

- Go 1.21 or later
- OCI account with appropriate permissions
- OCI API key configured
- A valid OCI compartment
- A subnet in your VCN
- (Optional) A capacity reservation

## Installation

1. Clone the repository:
```bash
git clone https://github.com/tomarkel/oci-insta-scale.git
cd oci-insta-scale
```

2. Install dependencies:
```bash
make deps
# or: go mod download
```

3. Configure your settings:
```bash
make setup-config
# or: cp config.example.yaml config.yaml
# Edit config.yaml with your OCI credentials and settings
```

4. Build the binaries:
```bash
make build
# or: go build -o oci-insta-scale main.go && go build -o capacity-manager capacity-manager.go
```

## Configuration

Edit `config.yaml` with your OCI settings:

### Required Fields:
- `tenancy_ocid`: Your OCI tenancy OCID
- `user_ocid`: Your user OCID
- `fingerprint`: Your API key fingerprint
- `private_key_path`: Path to your private key file
- `region`: OCI region (e.g., us-phoenix-1)
- `compartment_id`: Target compartment OCID
- `instance_settings.image_id`: OS image OCID
- `instance_settings.subnet_id`: Subnet OCID
- `instance_settings.availability_domain`: Availability domain name

### Optional Fields:
- `instance_settings.capacity_reservation_id`: Use a specific capacity reservation
- `instance_settings.fault_domain`: Specify fault domain
- `instance_settings.count`: Number of instances to create (default: 3)

### Finding Required OCIDs:

**List available images:**
```bash
oci compute image list --compartment-id <compartment-id> --output table
```

**List availability domains:**
```bash
oci iam availability-domain list --compartment-id <compartment-id>
```

**List capacity reservations:**
```bash
oci compute capacity-reservation list --compartment-id <compartment-id>
```

## Usage

### Using Make (Recommended):
```bash
# Build all binaries
make build

# Test configuration (dry-run)
make dry-run

# Create instances
make run

# View all available commands
make help
```

### Direct Usage:
```bash
# Basic usage
./oci-insta-scale

# With custom config file
./oci-insta-scale -config my-config.yaml

# Dry-run mode (preview without creating)
./oci-insta-scale -dry-run
```

## Capacity Reservations

To use a capacity reservation:

1. Create a capacity reservation in OCI Console or via the included tool:
```bash
# List existing reservations
./capacity-manager -list

# Create a new reservation
./capacity-manager -create \
  -name "my-reservation" \
  -ad "rgiR:US-ASHBURN-AD-1" \
  -shape "VM.Standard.E4.Flex" \
  -count 10 \
  -ocpus 1 \
  -memory 6

# Delete a reservation
./capacity-manager -delete -id "ocid1.capacityreservation.oc1.iad.xxxxx"
```

2. Add the capacity reservation OCID to your config:
```yaml
instance_settings:
  capacity_reservation_id: "ocid1.capacityreservation.oc1.iad.xxxxx"
```

3. Run the program - instances will be created using the reserved capacity.

## Example Output

```
Creating 3 instances...
[1/3] Creating instance: instance-1
  ✓ Created instance: instance-1 (OCID: ocid1.instance.oc1.iad.xxxxx)
[2/3] Creating instance: instance-2
  ✓ Created instance: instance-2 (OCID: ocid1.instance.oc1.iad.xxxxx)
[3/3] Creating instance: instance-3
  ✓ Created instance: instance-3 (OCID: ocid1.instance.oc1.iad.xxxxx)

=== Summary ===
Successfully created 3/3 instances

Waiting for instances to reach RUNNING state...
  Waiting for instance ocid1.instance.oc1.iad.xxxxx...
    Current state: PROVISIONING
    Current state: STARTING
    Current state: RUNNING
  ✓ Instance ocid1.instance.oc1.iad.xxxxx is RUNNING

All instances have been processed!
```

## Error Handling

The program will:
- Continue creating remaining instances if one fails
- Log errors for failed instances
- Provide a summary of successful vs failed creations
- Validate configuration before starting

## Cleanup

To terminate all created instances, you can use the OCI CLI:

```bash
# List instances with your prefix
oci compute instance list --compartment-id <compartment-id> \
  --display-name "instance-*" --query 'data[].id'

# Terminate a specific instance
oci compute instance terminate --instance-id <instance-id>
```

## License

MIT License

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

## Support

For issues related to:
- **This tool**: Open a GitHub issue
- **OCI API/SDK**: Check [OCI SDK documentation](https://docs.oracle.com/en-us/iaas/tools/go/latest/)
- **OCI Services**: Contact Oracle Cloud Support
