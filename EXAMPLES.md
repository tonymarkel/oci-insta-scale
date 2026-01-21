# Complete Configuration Examples

## Example 1: Basic Configuration (Minimal)

```yaml
# Minimal working configuration
tenancy_ocid: "ocid1.tenancy.oc1..your-tenancy-id"
user_ocid: "ocid1.user.oc1..your-user-id"
fingerprint: "aa:bb:cc:dd:ee:ff:11:22:33:44:55:66:77:88:99:00"
private_key_path: "/Users/you/.oci/oci_api_key.pem"
region: "us-phoenix-1"
compartment_id: "ocid1.compartment.oc1..your-compartment-id"

instance_settings:
  display_name_prefix: "instance"
  count: 3
  shape: "VM.Standard.E4.Flex"
  shape_config:
    ocpus: 1
    memory_in_gbs: 6
  image_id: "ocid1.image.oc1.phx.your-image-id"
  subnet_id: "ocid1.subnet.oc1.phx.your-subnet-id"
  availability_domain: "IYiP:PHX-AD-1"
  assign_public_ip: true
```

## Example 2: Production with Capacity Reservation

```yaml
# Production configuration with all features
tenancy_ocid: "ocid1.tenancy.oc1..your-tenancy-id"
user_ocid: "ocid1.user.oc1..your-user-id"
fingerprint: "aa:bb:cc:dd:ee:ff:11:22:33:44:55:66:77:88:99:00"
private_key_path: "/Users/you/.oci/oci_api_key.pem"
region: "us-ashburn-1"
compartment_id: "ocid1.compartment.oc1..your-compartment-id"

instance_settings:
  display_name_prefix: "prod-app"
  count: 20  # Create 20 instances
  
  # Use powerful flexible shape
  shape: "VM.Standard.E4.Flex"
  shape_config:
    ocpus: 4
    memory_in_gbs: 64
  
  # Oracle Linux 8 image
  image_id: "ocid1.image.oc1.iad.your-image-id"
  
  # Private subnet (no public IP)
  subnet_id: "ocid1.subnet.oc1.iad.your-private-subnet-id"
  assign_public_ip: false
  
  # Location
  availability_domain: "rgiR:US-ASHBURN-AD-1"
  fault_domain: "FAULT-DOMAIN-1"  # For HA
  
  # Use capacity reservation for guaranteed capacity
  capacity_reservation_id: "ocid1.capacityreservation.oc1.iad.your-reservation-id"
  
  # SSH access
  ssh_authorized_keys: |
    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAA... user@host
    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAA... admin@host
  
  # Custom metadata for cloud-init
  metadata:
    application: "web-server"
    environment: "production"
    team: "platform"
    version: "2.0"
  
  # Cost tracking tags
  freeform_tags:
    Project: "WebApp"
    CostCenter: "Engineering"
    Environment: "Production"
    Owner: "ops-team"
    ManagedBy: "oci-insta-scale"
```

## Example 3: Development/Testing Environment

```yaml
# Dev/test configuration - lower resources, temporary
tenancy_ocid: "ocid1.tenancy.oc1..your-tenancy-id"
user_ocid: "ocid1.user.oc1..your-user-id"
fingerprint: "aa:bb:cc:dd:ee:ff:11:22:33:44:55:66:77:88:99:00"
private_key_path: "/Users/you/.oci/oci_api_key.pem"
region: "us-phoenix-1"
compartment_id: "ocid1.compartment.oc1..dev-compartment-id"

instance_settings:
  display_name_prefix: "dev-test"
  count: 5  # Small number for testing
  
  # Minimal resources for cost savings
  shape: "VM.Standard.E4.Flex"
  shape_config:
    ocpus: 1  # Minimum
    memory_in_gbs: 6  # Minimum
  
  image_id: "ocid1.image.oc1.phx.your-image-id"
  subnet_id: "ocid1.subnet.oc1.phx.your-public-subnet-id"
  assign_public_ip: true  # For SSH access
  
  availability_domain: "IYiP:PHX-AD-1"
  
  ssh_authorized_keys: |
    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAA... dev@laptop
  
  freeform_tags:
    Environment: "Development"
    Temporary: "true"
    AutoDelete: "24hours"
```

## Example 4: High-Performance Computing (HPC)

```yaml
# HPC cluster configuration
tenancy_ocid: "ocid1.tenancy.oc1..your-tenancy-id"
user_ocid: "ocid1.user.oc1..your-user-id"
fingerprint: "aa:bb:cc:dd:ee:ff:11:22:33:44:55:66:77:88:99:00"
private_key_path: "/Users/you/.oci/oci_api_key.pem"
region: "us-ashburn-1"
compartment_id: "ocid1.compartment.oc1..hpc-compartment-id"

instance_settings:
  display_name_prefix: "hpc-node"
  count: 50  # Large cluster
  
  # High-performance shape
  shape: "VM.Standard.E4.Flex"
  shape_config:
    ocpus: 8  # Max performance
    memory_in_gbs: 128
  
  image_id: "ocid1.image.oc1.iad.your-hpc-image-id"
  subnet_id: "ocid1.subnet.oc1.iad.your-hpc-subnet-id"
  assign_public_ip: false  # Private network only
  
  availability_domain: "rgiR:US-ASHBURN-AD-2"
  
  # Use capacity reservation for guaranteed capacity
  capacity_reservation_id: "ocid1.capacityreservation.oc1.iad.hpc-reservation"
  
  ssh_authorized_keys: |
    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAA... hpc-admin@cluster
  
  metadata:
    cluster_name: "research-cluster-01"
    scheduler: "slurm"
    network_type: "rdma"
  
  freeform_tags:
    Project: "Research"
    Workload: "HPC"
    Grant: "NSF-12345"
```

## Example 5: ARM-based Instances (Ampere)

```yaml
# Ampere ARM-based instances for cost efficiency
tenancy_ocid: "ocid1.tenancy.oc1..your-tenancy-id"
user_ocid: "ocid1.user.oc1..your-user-id"
fingerprint: "aa:bb:cc:dd:ee:ff:11:22:33:44:55:66:77:88:99:00"
private_key_path: "/Users/you/.oci/oci_api_key.pem"
region: "us-phoenix-1"
compartment_id: "ocid1.compartment.oc1..your-compartment-id"

instance_settings:
  display_name_prefix: "arm-worker"
  count: 10
  
  # Ampere Altra ARM processors
  shape: "VM.Standard.A1.Flex"
  shape_config:
    ocpus: 4
    memory_in_gbs: 24
  
  # ARM-compatible image (Ubuntu for ARM or Oracle Linux ARM)
  image_id: "ocid1.image.oc1.phx.your-arm-image-id"
  subnet_id: "ocid1.subnet.oc1.phx.your-subnet-id"
  assign_public_ip: true
  
  availability_domain: "IYiP:PHX-AD-3"
  
  ssh_authorized_keys: |
    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAA... ops@control
  
  freeform_tags:
    Architecture: "ARM64"
    CostOptimized: "true"
```

## Example 6: Multi-Region Deployment (Separate Configs)

### us-phoenix.yaml
```yaml
tenancy_ocid: "ocid1.tenancy.oc1..your-tenancy-id"
user_ocid: "ocid1.user.oc1..your-user-id"
fingerprint: "aa:bb:cc:dd:ee:ff:11:22:33:44:55:66:77:88:99:00"
private_key_path: "/Users/you/.oci/oci_api_key.pem"
region: "us-phoenix-1"
compartment_id: "ocid1.compartment.oc1..your-compartment-id"

instance_settings:
  display_name_prefix: "phx-app"
  count: 10
  shape: "VM.Standard.E4.Flex"
  shape_config:
    ocpus: 2
    memory_in_gbs: 16
  image_id: "ocid1.image.oc1.phx.your-image-id"
  subnet_id: "ocid1.subnet.oc1.phx.your-subnet-id"
  availability_domain: "IYiP:PHX-AD-1"
  assign_public_ip: false
  
  freeform_tags:
    Region: "us-phoenix-1"
    Geography: "US-West"
```

### us-ashburn.yaml
```yaml
tenancy_ocid: "ocid1.tenancy.oc1..your-tenancy-id"
user_ocid: "ocid1.user.oc1..your-user-id"
fingerprint: "aa:bb:cc:dd:ee:ff:11:22:33:44:55:66:77:88:99:00"
private_key_path: "/Users/you/.oci/oci_api_key.pem"
region: "us-ashburn-1"
compartment_id: "ocid1.compartment.oc1..your-compartment-id"

instance_settings:
  display_name_prefix: "iad-app"
  count: 10
  shape: "VM.Standard.E4.Flex"
  shape_config:
    ocpus: 2
    memory_in_gbs: 16
  image_id: "ocid1.image.oc1.iad.your-image-id"
  subnet_id: "ocid1.subnet.oc1.iad.your-subnet-id"
  availability_domain: "rgiR:US-ASHBURN-AD-1"
  assign_public_ip: false
  
  freeform_tags:
    Region: "us-ashburn-1"
    Geography: "US-East"
```

### Deploy to both regions:
```bash
./oci-insta-scale -config us-phoenix.yaml
./oci-insta-scale -config us-ashburn.yaml
```

## Shape Quick Reference

### Flexible Shapes (Configurable OCPUs & Memory)
- **VM.Standard.E4.Flex**: Intel (1-64 OCPUs, 1-1024 GB)
- **VM.Standard.E5.Flex**: Intel (1-94 OCPUs, 1-1536 GB)
- **VM.Standard3.Flex**: AMD (1-32 OCPUs, 1-512 GB)
- **VM.Standard.A1.Flex**: Ampere ARM (1-80 OCPUs, 1-512 GB)

### Fixed Shapes (Common)
- **VM.Standard2.1**: 1 OCPU, 15 GB RAM
- **VM.Standard2.2**: 2 OCPUs, 30 GB RAM
- **VM.Standard2.4**: 4 OCPUs, 60 GB RAM
- **VM.Standard2.8**: 8 OCPUs, 120 GB RAM

## How to Find Your Values

### 1. Get your compartment ID
```bash
oci iam compartment list --all
```

### 2. Get availability domains
```bash
oci iam availability-domain list \
  --compartment-id <your-compartment-id>
```

### 3. Get image OCIDs
```bash
# Oracle Linux 8
oci compute image list \
  --compartment-id <compartment-id> \
  --operating-system "Oracle Linux" \
  --operating-system-version "8" \
  --shape "VM.Standard.E4.Flex"

# Ubuntu 22.04
oci compute image list \
  --compartment-id <compartment-id> \
  --operating-system "Canonical Ubuntu" \
  --operating-system-version "22.04"
```

### 4. Get subnet ID
```bash
oci network subnet list \
  --compartment-id <compartment-id> \
  --display-name "your-subnet-name"
```

### 5. Get or create capacity reservation
```bash
# List existing
./capacity-manager -list

# Create new
./capacity-manager -create \
  -name "my-reservation" \
  -ad "<availability-domain>" \
  -shape "VM.Standard.E4.Flex" \
  -count 20 \
  -ocpus 2 \
  -memory 16
```
