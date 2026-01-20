# OCI Instance Scaler

A Go program that provisions multiple OCI compute instances in parallel and can terminate them in bulk.

## Prerequisites

1. **OCI Account**: You need an active OCI account
2. **OCI CLI Configuration**: Set up OCI CLI credentials at `~/.oci/config`
3. **Go 1.21+**: Ensure Go is installed

## Setup

Install dependencies:
```bash
go mod download
```

## Building

```bash
go build -o oci-insta-scale
```

## Usage

### Creating Instances

```bash
./oci-insta-scale \
  -instances 5 \
  -name my-instance \
  -compartment <COMPARTMENT_ID> \
  -subnet <SUBNET_ID> \
  -image <IMAGE_ID> \
  -ad <AVAILABILITY_DOMAIN> \
  -shape VM.Standard.E4.Flex \
  -output instances.txt
```

#### Flags for Creation

- `-instances` (int): Number of instances to create (default: 1)
- `-name` (string): Base name for instances (default: "oci-instance")
- `-compartment` (string, required): OCI Compartment ID
- `-subnet` (string, required): OCI Subnet ID
- `-image` (string, required): OCI Image ID
- `-ad` (string, required): Availability Domain (e.g., `iad-ad-1` or `rgiR:US-ASHBURN-AD-2`)
- `-shape` (string): Instance shape (default: "VM.Standard.E4.Flex")
- `-output` (string): Output file for instance OCIDs (default: "instances.txt")

#### Example

```bash
./oci-insta-scale \
  -instances 10 \
  -name web-server \
  -compartment ocid1.compartment.oc1..example \
  -subnet ocid1.subnet.oc1..example \
  -image ocid1.image.oc1..example \
  -ad iad-ad-1 \
  -output my-instances.txt
```

The program will create instances and write their OCIDs to the specified output file (default: `instances.txt`).

### Terminating Instances

Terminate all instances listed in a file:

```bash
./oci-insta-scale terminate \
  -file instances.txt \
  -compartment <COMPARTMENT_ID> \
  -parallel 10
```

#### Flags for Termination

- `-file` (string): File containing instance OCIDs, one per line (default: "instances.txt")
- `-compartment` (string, required): OCI Compartment ID
- `-parallel` (int): Number of parallel termination operations (default: 10)

#### Example

```bash
./oci-insta-scale terminate \
  -file my-instances.txt \
  -compartment ocid1.compartment.oc1..example \
  -parallel 20
```

## How It Works

- **Parallel Execution**: Uses goroutines and `sync.WaitGroup` for concurrent operations
- **OCI SDK**: Uses the official OCI Go SDK v65
- **Configuration**: Reads OCI credentials from standard `~/.oci/config`
- **Instance Tracking**: Saves instance OCIDs to a file for later reference
- **Bulk Termination**: Terminate multiple instances concurrently with configurable parallelism
- **Results Tracking**: Displays real-time progress and summary

## Instance File Format

The output file contains one instance OCID per line:

```
ocid1.instance.oc1.iad.example1
ocid1.instance.oc1.iad.example2
ocid1.instance.oc1.iad.example3
```

You can edit this file to remove instances you want to keep before running the terminate command.

## Notes

- Ensure your OCI credentials are properly configured
- The parallelism level can be adjusted for both creation and termination
- Default shape configuration: 1 OCPU, 8GB memory (adjustable in code)
- Instance files support comments (lines starting with #) which are ignored during termination
