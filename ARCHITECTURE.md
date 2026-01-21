# Project Structure

```
oci-insta-scale/
├── main.go                    # Main instance creation program
├── capacity-manager.go        # Capacity reservation management utility
├── config.yaml                # Your configuration (gitignored)
├── config.example.yaml        # Example configuration template
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
├── README.md                  # Comprehensive documentation
├── QUICKSTART.md              # Quick start guide
├── manage-instances.sh        # Bash script for batch operations
└── images/                    # Screenshots/diagrams

# Binaries (after building)
├── oci-insta-scale           # Instance creation binary
└── capacity-manager          # Capacity manager binary
```

## File Descriptions

### Core Files

**main.go**
- Main application for creating instances
- Supports capacity reservations
- Configurable via YAML
- Includes dry-run mode
- Monitors instance status

**capacity-manager.go**
- Manage capacity reservations
- List, create, and delete reservations
- Shows reservation usage
- Supports flexible shapes

**config.yaml**
- Your actual configuration
- Contains OCI credentials
- Instance settings
- **Important**: Never commit this file (it's in .gitignore)

**config.example.yaml**
- Template configuration
- Safe to commit
- Shows all available options
- Copy to config.yaml and customize

### Documentation

**README.md**
- Complete feature documentation
- Installation instructions
- Configuration guide
- Usage examples
- Troubleshooting

**QUICKSTART.md**
- Step-by-step setup guide
- Common configuration examples
- Quick reference
- Shape table

### Utilities

**manage-instances.sh**
- Bash script for batch operations
- List, start, stop, terminate instances
- Uses OCI CLI
- Supports filtering by prefix
- Safety confirmations for destructive operations

## Workflow

### Initial Setup
1. Clone repository
2. Run `go mod download`
3. Copy `config.example.yaml` to `config.yaml`
4. Edit `config.yaml` with your OCI details
5. Build binaries: `go build`

### Creating Instances

**Without Capacity Reservation:**
```bash
# Test configuration
./oci-insta-scale -dry-run

# Create instances
./oci-insta-scale
```

**With Capacity Reservation:**
```bash
# Create reservation
./capacity-manager -create \
  -name "my-reservation" \
  -ad "rgiR:US-ASHBURN-AD-1" \
  -shape "VM.Standard.E4.Flex" \
  -count 10

# Add reservation OCID to config.yaml
# Then create instances
./oci-insta-scale
```

### Managing Instances
```bash
# List all instances
./manage-instances.sh list -c <compartment-id>

# Stop instances with prefix "test-"
./manage-instances.sh stop -c <compartment-id> -p "test-"

# Start instances
./manage-instances.sh start -c <compartment-id> -p "test-"

# Check status
./manage-instances.sh status -c <compartment-id>

# Terminate (requires confirmation)
./manage-instances.sh terminate -c <compartment-id> -p "test-"
```

### Cleanup
```bash
# Delete instances
./manage-instances.sh terminate -c <compartment-id> -y

# Delete capacity reservation
./capacity-manager -delete -id <reservation-ocid>
```

## Key Features

### Main Program (oci-insta-scale)
- ✓ Create arbitrary number of instances
- ✓ Use capacity reservations
- ✓ Support flexible shapes
- ✓ Dry-run mode
- ✓ Progress tracking
- ✓ Error handling and retry
- ✓ Wait for instance ready state
- ✓ Custom metadata and tags
- ✓ SSH key injection

### Capacity Manager
- ✓ List all reservations
- ✓ Create new reservations
- ✓ Delete reservations
- ✓ Show usage statistics
- ✓ Support flexible shapes

### Management Script
- ✓ Batch operations
- ✓ Filter by name prefix
- ✓ Start/stop instances
- ✓ Terminate with confirmation
- ✓ Status monitoring
- ✓ Color-coded output

## Security Notes

- **config.yaml is in .gitignore** - never commit credentials
- Use OCI IAM policies to limit permissions
- Store private keys securely (outside repository)
- Use separate compartments for different environments
- Apply appropriate security lists to subnets

## Integration Examples

### CI/CD Pipeline
```yaml
# GitHub Actions example
- name: Create test instances
  run: |
    ./oci-insta-scale -config ci-config.yaml
    
- name: Run tests
  run: ./run-tests.sh

- name: Cleanup
  run: |
    ./manage-instances.sh terminate \
      -c $COMPARTMENT_ID \
      -p "ci-test-" -y
```

### Terraform Integration
Use this tool to quickly provision instances, then manage with Terraform:
```bash
# Create instances
./oci-insta-scale

# Import to Terraform
terraform import oci_core_instance.instance[0] <instance-ocid>
```

## Best Practices

1. **Always test with dry-run first**
   ```bash
   ./oci-insta-scale -dry-run
   ```

2. **Use capacity reservations for guaranteed capacity**
   ```bash
   ./capacity-manager -create ...
   ```

3. **Tag your resources for cost tracking**
   ```yaml
   freeform_tags:
     Project: "MyProject"
     CostCenter: "Engineering"
   ```

4. **Use separate configs for different environments**
   ```bash
   ./oci-insta-scale -config production.yaml
   ./oci-insta-scale -config development.yaml
   ```

5. **Clean up unused resources**
   ```bash
   ./manage-instances.sh terminate -c <id> -p "old-"
   ```

## Support

- GitHub Issues: For bugs and feature requests
- OCI Documentation: https://docs.oracle.com/en-us/iaas/
- OCI SDK Go: https://github.com/oracle/oci-go-sdk
