# 🚀 OCI Insta-Scale - Complete Go Program

A production-ready Go application for creating and managing Oracle Cloud Infrastructure (OCI) compute instances at scale using capacity reservations.

## ✨ What's Been Created

### Core Applications (Go)

1. **main.go** - Main instance creation program
   - Creates arbitrary number of instances
   - Supports OCI capacity reservations
   - Configurable via YAML
   - Dry-run mode for testing
   - Progress tracking and monitoring
   - Error handling and recovery

2. **capacity-manager.go** - Capacity reservation utility
   - List all reservations
   - Create new reservations
   - Delete reservations
   - Show usage statistics
   - Support for flexible shapes

### Utilities & Scripts

3. **manage-instances.sh** - Bash script for batch operations
   - List, start, stop, terminate instances
   - Filter by name prefix
   - Safety confirmations
   - Color-coded output

4. **Makefile** - Build automation
   - One-command builds
   - Common task shortcuts
   - Development helpers

### Configuration & Documentation

5. **config.yaml** - Your actual configuration
   - OCI credentials
   - Instance settings
   - Capacity reservation settings
   - *(gitignored for security)*

6. **config.example.yaml** - Template configuration
   - All available options
   - Comments and examples
   - Safe to commit

7. **Documentation Files**:
   - `README.md` - Complete feature documentation
   - `QUICKSTART.md` - Step-by-step setup guide
   - `ARCHITECTURE.md` - Project structure and patterns
   - `WORKFLOW.md` - Visual workflow diagrams

## 📦 Project Structure

```
oci-insta-scale/
├── 🔧 Core Go Files
│   ├── main.go                 # Instance creation
│   ├── capacity-manager.go     # Reservation manager
│   ├── go.mod                  # Dependencies
│   └── go.sum                  # Checksums
│
├── 🔨 Build Artifacts
│   ├── oci-insta-scale         # Main binary
│   └── capacity-manager        # Manager binary
│
├── ⚙️  Configuration
│   ├── config.yaml             # Your config (gitignored)
│   └── config.example.yaml     # Template
│
├── 📜 Scripts
│   ├── manage-instances.sh     # Batch operations
│   └── Makefile                # Build automation
│
└── 📚 Documentation
    ├── README.md               # Main docs
    ├── QUICKSTART.md          # Quick setup
    ├── ARCHITECTURE.md        # Structure
    └── WORKFLOW.md            # Diagrams
```

## 🎯 Key Features

### Instance Creation
✅ Create 1-1000+ instances with one command  
✅ Use capacity reservations for guaranteed capacity  
✅ Flexible shape configuration (OCPUs, memory)  
✅ Support for SSH keys, metadata, and tags  
✅ Dry-run mode for testing  
✅ Progress tracking with status updates  
✅ Error handling with detailed logging  
✅ Auto-wait for RUNNING state  

### Capacity Management
✅ List all capacity reservations  
✅ Create new reservations  
✅ Delete unused reservations  
✅ View usage statistics  
✅ Support for flexible and fixed shapes  

### Operations
✅ Batch start/stop/terminate  
✅ Filter by name prefix  
✅ Status monitoring  
✅ Safety confirmations  

## 🚀 Quick Start

### 1. Setup (One Time)
```bash
# Clone and setup
git clone <repo>
cd oci-insta-scale

# Install dependencies
make deps

# Create configuration
make setup-config
vi config.yaml  # Add your OCI credentials

# Build binaries
make build
```

### 2. Create Instances

**Option A: Without Capacity Reservation**
```bash
# Test first
make dry-run

# Create instances
make run
```

**Option B: With Capacity Reservation**
```bash
# Create reservation
./capacity-manager -create \
  -name "prod-reservation" \
  -ad "rgiR:US-ASHBURN-AD-1" \
  -shape "VM.Standard.E4.Flex" \
  -count 10 \
  -ocpus 1 \
  -memory 6

# Add reservation OCID to config.yaml
# Then create instances
make run
```

### 3. Manage Instances
```bash
# List instances
./manage-instances.sh list -c <compartment-id>

# Check status
./manage-instances.sh status -c <compartment-id>

# Stop instances
./manage-instances.sh stop -c <compartment-id> -p "test-"

# Terminate (with confirmation)
./manage-instances.sh terminate -c <compartment-id> -p "test-"
```

## 📋 Configuration Example

```yaml
# config.yaml
tenancy_ocid: "ocid1.tenancy.oc1..aaaa..."
user_ocid: "ocid1.user.oc1..aaaa..."
fingerprint: "4a:df:0b:63:f9:f4:ae:52:..."
private_key_path: "/Users/you/.oci/key.pem"
region: "us-phoenix-1"
compartment_id: "ocid1.compartment.oc1..aaaa..."

instance_settings:
  display_name_prefix: "my-app"
  count: 10
  shape: "VM.Standard.E4.Flex"
  shape_config:
    ocpus: 2
    memory_in_gbs: 12
  image_id: "ocid1.image.oc1.iad.aaaa..."
  subnet_id: "ocid1.subnet.oc1.iad.aaaa..."
  availability_domain: "rgiR:US-ASHBURN-AD-1"
  
  # Optional: Use capacity reservation
  capacity_reservation_id: "ocid1.capacityreservation..."
  
  # Optional: Tags
  freeform_tags:
    Project: "MyProject"
    Environment: "Production"
```

## 🔍 Example Output

```bash
$ ./oci-insta-scale

Creating 10 instances...
[1/10] Creating instance: my-app-1
  ✓ Created instance: my-app-1 (OCID: ocid1.instance...)
[2/10] Creating instance: my-app-2
  ✓ Created instance: my-app-2 (OCID: ocid1.instance...)
...
[10/10] Creating instance: my-app-10
  ✓ Created instance: my-app-10 (OCID: ocid1.instance...)

=== Summary ===
Successfully created 10/10 instances

Waiting for instances to reach RUNNING state...
  ✓ Instance my-app-1 is RUNNING
  ✓ Instance my-app-2 is RUNNING
  ...
  ✓ Instance my-app-10 is RUNNING

All instances have been processed!
```

## 🛠️ Make Commands

```bash
make build              # Build all binaries
make clean              # Remove build artifacts
make deps               # Download dependencies
make run                # Create instances
make dry-run            # Test without creating
make list-reservations  # List capacity reservations
make setup-config       # Create config from template
make help               # Show all commands
```

## 📚 Documentation Guide

- **Start here**: [README.md](README.md) - Complete feature documentation
- **Quick setup**: [QUICKSTART.md](QUICKSTART.md) - Step-by-step instructions
- **Understanding**: [ARCHITECTURE.md](ARCHITECTURE.md) - Project structure
- **Workflows**: [WORKFLOW.md](WORKFLOW.md) - Visual diagrams

## 🔐 Security Notes

✅ `config.yaml` is in `.gitignore` - never committed  
✅ Use OCI IAM policies for least privilege  
✅ Store private keys securely  
✅ Separate configs for dev/staging/prod  
✅ Apply security lists to subnets  

## 🎓 What You Can Do Now

1. **Create instances at scale**
   - From 1 to 1000+ instances
   - With or without capacity reservations

2. **Manage capacity**
   - Create reservations for guaranteed capacity
   - List and monitor usage
   - Delete when no longer needed

3. **Batch operations**
   - Start/stop multiple instances
   - Terminate by name prefix
   - Monitor status across fleet

4. **Automate workflows**
   - CI/CD integration ready
   - Configuration management
   - Cost optimization

## 💡 Use Cases

- **Burst workloads**: Create 100+ instances for batch processing
- **Testing**: Spin up temporary test environments
- **Development**: Quick dev instance provisioning
- **Production**: Guaranteed capacity with reservations
- **CI/CD**: Automated test infrastructure
- **Research**: Large-scale computation clusters

## 🆘 Troubleshooting

**Build Issues**
```bash
make clean
make deps
make build
```

**Configuration Issues**
```bash
# Verify OCI CLI works
oci iam region list

# Test with dry-run
make dry-run
```

**Authentication Issues**
- Check private key path and permissions
- Verify fingerprint matches API key
- Ensure user has required policies

**Capacity Issues**
- Create a capacity reservation first
- Try different availability domain
- Use different shape if unavailable

## 📞 Support

- **Issues**: GitHub Issues
- **OCI Docs**: https://docs.oracle.com/en-us/iaas/
- **OCI Go SDK**: https://github.com/oracle/oci-go-sdk

## 🎉 You're Ready!

Everything is set up and ready to use. Just:

1. Edit `config.yaml` with your OCI credentials
2. Run `make dry-run` to test
3. Run `make run` to create instances

Enjoy scaling on OCI! 🚀
