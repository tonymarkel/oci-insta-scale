# 🎯 Setup Checklist

Use this checklist to get started with OCI Insta-Scale.

## ☑️ Prerequisites

- [ ] Go 1.21 or later installed
  ```bash
  go version  # Should show 1.21 or higher
  ```

- [ ] OCI CLI installed (optional, but helpful)
  ```bash
  oci --version
  ```

- [ ] OCI account with active subscription

- [ ] API keys configured for OCI
  - Private key file (.pem)
  - Fingerprint from OCI Console
  - User OCID
  - Tenancy OCID

## ☑️ OCI Account Setup

- [ ] Created compartment for your instances
  ```
  OCI Console → Identity → Compartments → Create Compartment
  ```

- [ ] Noted compartment OCID
  ```
  Copy from compartment details page
  ```

- [ ] Created VCN (Virtual Cloud Network)
  ```
  OCI Console → Networking → Virtual Cloud Networks → Start VCN Wizard
  ```

- [ ] Created subnet in VCN
  ```
  Note the subnet OCID
  ```

- [ ] Configured security lists/network security groups
  ```
  Allow SSH (port 22) and required application ports
  ```

- [ ] Found availability domain names
  ```bash
  oci iam availability-domain list --compartment-id <id>
  ```

- [ ] Found image OCID for desired OS
  ```bash
  oci compute image list --compartment-id <id> \
    --operating-system "Oracle Linux" \
    --operating-system-version "8" --output table
  ```

## ☑️ Project Setup

- [ ] Cloned repository
  ```bash
  git clone <repo>
  cd oci-insta-scale
  ```

- [ ] Downloaded dependencies
  ```bash
  make deps
  # or: go mod download
  ```

- [ ] Created config.yaml from template
  ```bash
  make setup-config
  # or: cp config.example.yaml config.yaml
  ```

- [ ] Edited config.yaml with your values:
  - [ ] tenancy_ocid
  - [ ] user_ocid
  - [ ] fingerprint
  - [ ] private_key_path
  - [ ] region
  - [ ] compartment_id
  - [ ] image_id
  - [ ] subnet_id
  - [ ] availability_domain

- [ ] Built binaries
  ```bash
  make build
  ```

## ☑️ Verify Configuration

- [ ] Private key file exists and is readable
  ```bash
  ls -l ~/.oci/*.pem
  # Should show: -rw------- (permissions 600)
  ```

- [ ] OCI credentials work (if CLI installed)
  ```bash
  oci iam region list
  # Should show list of regions
  ```

- [ ] Test with dry-run
  ```bash
  make dry-run
  # Should show configuration without errors
  ```

## ☑️ Optional: Capacity Reservation

If you want guaranteed capacity:

- [ ] Decide on reservation size and shape

- [ ] Create capacity reservation
  ```bash
  ./capacity-manager -create \
    -name "my-reservation" \
    -ad "<availability-domain>" \
    -shape "VM.Standard.E4.Flex" \
    -count 10 \
    -ocpus 2 \
    -memory 16
  ```

- [ ] Added reservation OCID to config.yaml
  ```yaml
  instance_settings:
    capacity_reservation_id: "ocid1.capacityreservation..."
  ```

## ☑️ First Run

- [ ] Start with small count (e.g., 2-3 instances)
  ```yaml
  instance_settings:
    count: 2
  ```

- [ ] Run dry-run to verify
  ```bash
  make dry-run
  ```

- [ ] Create instances
  ```bash
  make run
  ```

- [ ] Verify instances created
  ```bash
  # Using OCI CLI
  oci compute instance list --compartment-id <id> \
    --display-name "instance-*"
  
  # Or using management script
  ./manage-instances.sh list -c <compartment-id>
  ```

- [ ] Test SSH access to one instance
  ```bash
  ssh opc@<public-ip>
  ```

- [ ] Clean up test instances
  ```bash
  ./manage-instances.sh terminate -c <compartment-id> -y
  ```

## ☑️ Production Readiness

- [ ] Tested with small count successfully

- [ ] Decided on final configuration:
  - [ ] Instance count
  - [ ] Shape and resources
  - [ ] Network settings
  - [ ] Tags for cost tracking

- [ ] Created capacity reservation (if needed)

- [ ] Set up monitoring/alerting in OCI Console

- [ ] Documented your deployment process

- [ ] Configured backup configs for different environments
  ```bash
  config.yaml          # Production
  dev-config.yaml      # Development
  staging-config.yaml  # Staging
  ```

## ☑️ Post-Deployment

- [ ] Verified all instances are RUNNING

- [ ] Tested application connectivity

- [ ] Configured load balancer (if applicable)

- [ ] Set up monitoring and logging

- [ ] Documented instance details

- [ ] Added cost tracking tags

## 🚨 Common Issues Checklist

If something doesn't work, check:

- [ ] Private key path is correct and file exists
- [ ] Private key has correct permissions (600)
- [ ] Fingerprint matches the key in OCI Console
- [ ] User OCID is correct
- [ ] Tenancy OCID is correct
- [ ] Region is correct
- [ ] Compartment OCID is correct and exists
- [ ] Image OCID is valid for the region
- [ ] Subnet OCID is valid and in the correct AD
- [ ] Availability domain name is correct
- [ ] Security lists allow required traffic
- [ ] IAM policies grant required permissions
- [ ] Service limits allow the instance count
- [ ] Capacity is available in the AD

## 📋 Quick Reference

### Find OCIDs
```bash
# Compartments
oci iam compartment list --all

# Availability Domains
oci iam availability-domain list --compartment-id <id>

# Images
oci compute image list --compartment-id <id>

# Subnets
oci network subnet list --compartment-id <id>

# Instances
oci compute instance list --compartment-id <id>
```

### Common Commands
```bash
# Build
make build

# Test
make dry-run

# Create instances
make run

# List reservations
make list-reservations

# Manage instances
./manage-instances.sh list -c <compartment-id>
./manage-instances.sh status -c <compartment-id>
./manage-instances.sh stop -c <compartment-id> -p "prefix-"
./manage-instances.sh start -c <compartment-id> -p "prefix-"
./manage-instances.sh terminate -c <compartment-id> -p "prefix-"
```

## ✅ Success Criteria

You're ready for production when:

- [x] Dry-run shows correct configuration
- [x] Test instances create successfully
- [x] You can SSH into created instances
- [x] Network connectivity works as expected
- [x] Cost tracking tags are applied
- [x] Monitoring is configured
- [x] Cleanup process is documented

## 📚 Next Steps

After successful setup:

1. **Scale Up**: Increase instance count in config
2. **Automate**: Integrate with CI/CD pipelines
3. **Monitor**: Set up CloudWatch/OCI Monitoring
4. **Optimize**: Review costs and adjust resources
5. **Document**: Keep deployment notes updated

## 🆘 Need Help?

- Review [README.md](README.md) for detailed docs
- Check [QUICKSTART.md](QUICKSTART.md) for step-by-step guide
- See [EXAMPLES.md](EXAMPLES.md) for configuration examples
- Read [WORKFLOW.md](WORKFLOW.md) for architecture details
- Open GitHub issue for bugs or questions

---

**Print this checklist and mark items as you complete them!**
