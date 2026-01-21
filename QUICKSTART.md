# Quick Start Guide

## 1. Initial Setup

```bash
# Install dependencies
go mod download

# Copy and configure
cp config.example.yaml config.yaml
# Edit config.yaml with your OCI credentials
```

## 2. Find Required Values

### Get Availability Domains
```bash
oci iam availability-domain list --compartment-id <compartment-id>
```

### Get Image OCIDs (Oracle Linux 8)
```bash
oci compute image list \
  --compartment-id <compartment-id> \
  --operating-system "Oracle Linux" \
  --operating-system-version "8" \
  --output table
```

### Get Subnet OCID
```bash
oci network subnet list --compartment-id <compartment-id> --output table
```

## 3. Build the Tools

```bash
# Build main instance creator
go build -o oci-insta-scale main.go

# Build capacity reservation manager
go build -o capacity-manager capacity-manager.go
```

## 4. Manage Capacity Reservations (Optional)

```bash
# List existing reservations
./capacity-manager -list

# Create a new reservation for 10 instances
./capacity-manager -create \
  -name "prod-reservation" \
  -ad "rgiR:US-ASHBURN-AD-1" \
  -shape "VM.Standard.E4.Flex" \
  -count 10 \
  -ocpus 1 \
  -memory 6

# The tool will output the reservation OCID - add it to config.yaml
```

## 5. Create Instances

```bash
# Test with dry-run first
./oci-insta-scale -dry-run

# Create instances
./oci-insta-scale

# Use custom config
./oci-insta-scale -config production.yaml
```

## 6. Verify Instances

```bash
# List instances
oci compute instance list \
  --compartment-id <compartment-id> \
  --display-name "instance-*" \
  --output table

# Get instance details
oci compute instance get --instance-id <instance-id>
```

## 7. Cleanup

```bash
# List instances to delete
oci compute instance list \
  --compartment-id <compartment-id> \
  --display-name "instance-*" \
  --query 'data[].id'

# Terminate an instance
oci compute instance terminate --instance-id <instance-id> --force

# Delete capacity reservation (when no longer needed)
./capacity-manager -delete -id <reservation-ocid>
```

## Configuration Examples

### Minimal Configuration (3 instances)
```yaml
instance_settings:
  display_name_prefix: "test"
  count: 3
  shape: "VM.Standard.E4.Flex"
  shape_config:
    ocpus: 1
    memory_in_gbs: 6
  image_id: "ocid1.image.oc1.iad.xxxxx"
  subnet_id: "ocid1.subnet.oc1.iad.xxxxx"
  availability_domain: "rgiR:US-ASHBURN-AD-1"
```

### With Capacity Reservation
```yaml
instance_settings:
  display_name_prefix: "prod"
  count: 10
  shape: "VM.Standard.E4.Flex"
  shape_config:
    ocpus: 2
    memory_in_gbs: 12
  image_id: "ocid1.image.oc1.iad.xxxxx"
  subnet_id: "ocid1.subnet.oc1.iad.xxxxx"
  availability_domain: "rgiR:US-ASHBURN-AD-1"
  capacity_reservation_id: "ocid1.capacityreservation.oc1.iad.xxxxx"
  freeform_tags:
    Environment: "production"
```

### High Availability (with Fault Domain)
```yaml
instance_settings:
  display_name_prefix: "ha-app"
  count: 5
  shape: "VM.Standard.E4.Flex"
  shape_config:
    ocpus: 4
    memory_in_gbs: 24
  image_id: "ocid1.image.oc1.iad.xxxxx"
  subnet_id: "ocid1.subnet.oc1.iad.xxxxx"
  availability_domain: "rgiR:US-ASHBURN-AD-1"
  fault_domain: "FAULT-DOMAIN-1"
  assign_public_ip: false
```

## Common Shapes

| Shape | Type | OCPUs | Memory (GB) | Flexible |
|-------|------|-------|-------------|----------|
| VM.Standard.E4.Flex | General Purpose | 1-64 | 1-1024 | Yes |
| VM.Standard.E5.Flex | General Purpose | 1-94 | 1-1536 | Yes |
| VM.Standard3.Flex | General Purpose | 1-32 | 1-512 | Yes |
| VM.Standard.A1.Flex | Ampere Arm | 1-80 | 1-512 | Yes |
| VM.Standard2.1 | General Purpose | 1 | 15 | No |
| VM.Standard2.2 | General Purpose | 2 | 30 | No |
| BM.Standard.E4.128 | Bare Metal | 128 | 2048 | No |

## Troubleshooting

### "Out of capacity" error
- Create a capacity reservation first
- Try a different availability domain
- Use a different shape

### "Subnet not found" error
- Verify subnet OCID in config
- Check compartment permissions
- Ensure subnet is in the correct availability domain

### "Invalid shape" error
- Verify shape name is correct
- Check if shape is available in your tenancy
- For flexible shapes, ensure OCPUs and memory are within limits

### Authentication errors
- Verify OCI CLI is configured: `oci setup config`
- Check private key path and permissions
- Verify fingerprint matches your API key
