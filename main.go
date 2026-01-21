package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"gopkg.in/yaml.v3"
)

type Config struct {
	TenancyOCID      string `yaml:"tenancy_ocid"`
	UserOCID         string `yaml:"user_ocid"`
	Fingerprint      string `yaml:"fingerprint"`
	PrivateKeyPath   string `yaml:"private_key_path"`
	Region           string `yaml:"region"`
	CompartmentID    string `yaml:"compartment_id"`
	InstanceSettings struct {
		DisplayNamePrefix string `yaml:"display_name_prefix"`
		Count             int    `yaml:"count"`
		Shape             string `yaml:"shape"`
		ShapeConfig       struct {
			OCPUs        float32 `yaml:"ocpus"`
			MemoryInGBs  float32 `yaml:"memory_in_gbs"`
		} `yaml:"shape_config"`
		ImageID                string            `yaml:"image_id"`
		SubnetID               string            `yaml:"subnet_id"`
		AssignPublicIP         bool              `yaml:"assign_public_ip"`
		SSHAuthorizedKeys      string            `yaml:"ssh_authorized_keys"`
		AvailabilityDomain     string            `yaml:"availability_domain"`
		CapacityReservationID  *string           `yaml:"capacity_reservation_id,omitempty"`
		FaultDomain            *string           `yaml:"fault_domain,omitempty"`
		Metadata               map[string]string `yaml:"metadata"`
		FreeformTags           map[string]string `yaml:"freeform_tags"`
	} `yaml:"instance_settings"`
}

func main() {
	configFile := flag.String("config", "config.yaml", "Path to configuration file")
	dryRun := flag.Bool("dry-run", false, "Perform a dry run without creating instances")
	flag.Parse()

	// Load configuration
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Read private key file
	privateKey, err := os.ReadFile(config.PrivateKeyPath)
	if err != nil {
		log.Fatalf("Failed to read private key file: %v", err)
	}

	// Create OCI compute client
	configProvider := common.NewRawConfigurationProvider(
		config.TenancyOCID,
		config.UserOCID,
		config.Region,
		config.Fingerprint,
		string(privateKey),
		nil,
	)

	computeClient, err := core.NewComputeClientWithConfigurationProvider(configProvider)
	if err != nil {
		log.Fatalf("Failed to create compute client: %v", err)
	}

	ctx := context.Background()

	if *dryRun {
		fmt.Println("=== DRY RUN MODE ===")
		fmt.Printf("Would create %d instances with the following configuration:\n", config.InstanceSettings.Count)
		fmt.Printf("  Shape: %s\n", config.InstanceSettings.Shape)
		fmt.Printf("  OCPUs: %.1f\n", config.InstanceSettings.ShapeConfig.OCPUs)
		fmt.Printf("  Memory: %.1f GB\n", config.InstanceSettings.ShapeConfig.MemoryInGBs)
		fmt.Printf("  Image: %s\n", config.InstanceSettings.ImageID)
		fmt.Printf("  Subnet: %s\n", config.InstanceSettings.SubnetID)
		fmt.Printf("  Availability Domain: %s\n", config.InstanceSettings.AvailabilityDomain)
		if config.InstanceSettings.CapacityReservationID != nil {
			fmt.Printf("  Capacity Reservation: %s\n", *config.InstanceSettings.CapacityReservationID)
		}
		return
	}

	// Create instances
	fmt.Printf("Creating %d instances...\n", config.InstanceSettings.Count)
	instances := make([]string, 0, config.InstanceSettings.Count)

	for i := 0; i < config.InstanceSettings.Count; i++ {
		displayName := fmt.Sprintf("%s-%d", config.InstanceSettings.DisplayNamePrefix, i+1)
		
		fmt.Printf("[%d/%d] Creating instance: %s\n", i+1, config.InstanceSettings.Count, displayName)
		
		instanceID, err := createInstance(ctx, computeClient, config, displayName)
		if err != nil {
			log.Printf("Failed to create instance %s: %v", displayName, err)
			continue
		}
		
		instances = append(instances, instanceID)
		fmt.Printf("  ✓ Created instance: %s (OCID: %s)\n", displayName, instanceID)
		
		// Small delay to avoid rate limiting
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Successfully created %d/%d instances\n", len(instances), config.InstanceSettings.Count)
	
	// Wait for instances to be running (optional)
	if len(instances) > 0 {
		fmt.Println("\nWaiting for instances to reach RUNNING state...")
		waitForInstances(ctx, computeClient, instances)
	}
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate required fields
	if config.CompartmentID == "" {
		return nil, fmt.Errorf("compartment_id is required")
	}
	if config.InstanceSettings.ImageID == "" {
		return nil, fmt.Errorf("instance_settings.image_id is required")
	}
	if config.InstanceSettings.SubnetID == "" {
		return nil, fmt.Errorf("instance_settings.subnet_id is required")
	}
	if config.InstanceSettings.AvailabilityDomain == "" {
		return nil, fmt.Errorf("instance_settings.availability_domain is required")
	}

	return &config, nil
}

func createInstance(ctx context.Context, client core.ComputeClient, config *Config, displayName string) (string, error) {
	// Build launch instance details
	launchDetails := core.LaunchInstanceDetails{
		AvailabilityDomain: common.String(config.InstanceSettings.AvailabilityDomain),
		CompartmentId:      common.String(config.CompartmentID),
		DisplayName:        common.String(displayName),
		Shape:              common.String(config.InstanceSettings.Shape),
		ShapeConfig: &core.LaunchInstanceShapeConfigDetails{
			Ocpus:       common.Float32(config.InstanceSettings.ShapeConfig.OCPUs),
			MemoryInGBs: common.Float32(config.InstanceSettings.ShapeConfig.MemoryInGBs),
		},
		CreateVnicDetails: &core.CreateVnicDetails{
			SubnetId:       common.String(config.InstanceSettings.SubnetID),
			AssignPublicIp: common.Bool(config.InstanceSettings.AssignPublicIP),
		},
	}

	// Set source details
	launchDetails.SourceDetails = core.InstanceSourceViaImageDetails{
		ImageId: common.String(config.InstanceSettings.ImageID),
	}

	// Initialize metadata map if needed
	if config.InstanceSettings.Metadata != nil {
		launchDetails.Metadata = config.InstanceSettings.Metadata
	}

	// Initialize tags if provided
	if config.InstanceSettings.FreeformTags != nil {
		launchDetails.FreeformTags = config.InstanceSettings.FreeformTags
	}

	// Add SSH keys if provided
	if config.InstanceSettings.SSHAuthorizedKeys != "" {
		if launchDetails.Metadata == nil {
			launchDetails.Metadata = make(map[string]string)
		}
		launchDetails.Metadata["ssh_authorized_keys"] = config.InstanceSettings.SSHAuthorizedKeys
	}

	// Add capacity reservation if specified
	if config.InstanceSettings.CapacityReservationID != nil {
		launchDetails.CapacityReservationId = config.InstanceSettings.CapacityReservationID
	}

	// Add fault domain if specified
	if config.InstanceSettings.FaultDomain != nil {
		launchDetails.FaultDomain = config.InstanceSettings.FaultDomain
	}

	request := core.LaunchInstanceRequest{
		LaunchInstanceDetails: launchDetails,
	}

	// Launch the instance
	response, err := client.LaunchInstance(ctx, request)
	if err != nil {
		return "", fmt.Errorf("failed to launch instance: %w", err)
	}

	return *response.Instance.Id, nil
}

func waitForInstances(ctx context.Context, client core.ComputeClient, instanceIDs []string) {
	for _, instanceID := range instanceIDs {
		fmt.Printf("  Waiting for instance %s...\n", instanceID)
		
		for {
			request := core.GetInstanceRequest{
				InstanceId: common.String(instanceID),
			}
			
			response, err := client.GetInstance(ctx, request)
			if err != nil {
				log.Printf("  Error checking instance status: %v", err)
				break
			}
			
			state := response.Instance.LifecycleState
			fmt.Printf("    Current state: %s\n", state)
			
			if state == core.InstanceLifecycleStateRunning {
				fmt.Printf("  ✓ Instance %s is RUNNING\n", instanceID)
				break
			} else if state == core.InstanceLifecycleStateTerminated || 
					   state == core.InstanceLifecycleStateTerminating {
				log.Printf("  ✗ Instance %s is in terminal state: %s\n", instanceID, state)
				break
			}
			
			time.Sleep(10 * time.Second)
		}
	}
	
	fmt.Println("\nAll instances have been processed!")
}
