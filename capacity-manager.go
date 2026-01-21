package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"gopkg.in/yaml.v3"
)

type ReservationConfig struct {
	TenancyOCID    string `yaml:"tenancy_ocid"`
	UserOCID       string `yaml:"user_ocid"`
	Fingerprint    string `yaml:"fingerprint"`
	PrivateKeyPath string `yaml:"private_key_path"`
	Region         string `yaml:"region"`
	CompartmentID  string `yaml:"compartment_id"`
}

func main() {
	configFile := flag.String("config", "config.yaml", "Path to configuration file")
	list := flag.Bool("list", false, "List all capacity reservations")
	create := flag.Bool("create", false, "Create a new capacity reservation")
	delete := flag.Bool("delete", false, "Delete a capacity reservation")
	reservationID := flag.String("id", "", "Capacity reservation OCID (for delete)")
	
	// Creation parameters
	displayName := flag.String("name", "my-capacity-reservation", "Display name for the reservation")
	availabilityDomain := flag.String("ad", "", "Availability domain")
	shape := flag.String("shape", "VM.Standard.E4.Flex", "Instance shape")
	count := flag.Int64("count", 10, "Number of instances to reserve")
	ocpus := flag.Float64("ocpus", 1, "OCPUs per instance (for flexible shapes)")
	memoryInGBs := flag.Float64("memory", 6, "Memory in GB per instance (for flexible shapes)")
	
	flag.Parse()

	// Load configuration
	config, err := loadReservationConfig(*configFile)
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

	// Execute command
	if *list {
		listReservations(ctx, computeClient, config.CompartmentID)
	} else if *create {
		if *availabilityDomain == "" {
			log.Fatal("Availability domain (-ad) is required for creation")
		}
		createReservation(ctx, computeClient, config.CompartmentID, *displayName, *availabilityDomain, *shape, *count, float32(*ocpus), float32(*memoryInGBs))
	} else if *delete {
		if *reservationID == "" {
			log.Fatal("Reservation ID (-id) is required for deletion")
		}
		deleteReservation(ctx, computeClient, *reservationID)
	} else {
		flag.Usage()
		os.Exit(1)
	}
}

func loadReservationConfig(path string) (*ReservationConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ReservationConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func listReservations(ctx context.Context, client core.ComputeClient, compartmentID string) {
	fmt.Println("Listing capacity reservations...")
	fmt.Println(strings.Repeat("=", 80))

	request := core.ListComputeCapacityReservationsRequest{
		CompartmentId: common.String(compartmentID),
	}

	response, err := client.ListComputeCapacityReservations(ctx, request)
	if err != nil {
		log.Fatalf("Failed to list capacity reservations: %v", err)
	}

	if len(response.Items) == 0 {
		fmt.Println("No capacity reservations found.")
		return
	}

	for _, reservation := range response.Items {
		fmt.Printf("\nName: %s\n", *reservation.DisplayName)
		fmt.Printf("  OCID: %s\n", *reservation.Id)
		fmt.Printf("  State: %s\n", reservation.LifecycleState)
		fmt.Printf("  Availability Domain: %s\n", *reservation.AvailabilityDomain)
		fmt.Printf("  Reserved Instances: %d\n", *reservation.ReservedInstanceCount)
		fmt.Printf("  Used Instances: %d\n", *reservation.UsedInstanceCount)
		
		if reservation.TimeCreated != nil {
			fmt.Printf("  Created: %s\n", reservation.TimeCreated.Format("2006-01-02 15:04:05"))
		}
	}
	fmt.Println(strings.Repeat("=", 80))
}

func createReservation(ctx context.Context, client core.ComputeClient, compartmentID, displayName, availabilityDomain, shape string, count int64, ocpus, memoryInGBs float32) {
	fmt.Printf("Creating capacity reservation: %s\n", displayName)

	// Build instance reservation config
	instanceConfig := core.InstanceReservationConfigDetails{
		InstanceShape: common.String(shape),
		ReservedCount: common.Int64(count),
	}

	// Add shape config for flexible shapes
	if isFlexibleShape(shape) {
		instanceConfig.InstanceShapeConfig = &core.InstanceReservationShapeConfigDetails{
			Ocpus:       common.Float32(ocpus),
			MemoryInGBs: common.Float32(memoryInGBs),
		}
	}

	request := core.CreateComputeCapacityReservationRequest{
		CreateComputeCapacityReservationDetails: core.CreateComputeCapacityReservationDetails{
			CompartmentId:              common.String(compartmentID),
			DisplayName:                common.String(displayName),
			AvailabilityDomain:         common.String(availabilityDomain),
			InstanceReservationConfigs: []core.InstanceReservationConfigDetails{instanceConfig},
		},
	}

	response, err := client.CreateComputeCapacityReservation(ctx, request)
	if err != nil {
		log.Fatalf("Failed to create capacity reservation: %v", err)
	}

	fmt.Println("✓ Capacity reservation created successfully!")
	fmt.Printf("  OCID: %s\n", *response.Id)
	fmt.Printf("  State: %s\n", response.LifecycleState)
	fmt.Printf("\nAdd this to your config.yaml:")
	fmt.Printf("  capacity_reservation_id: \"%s\"\n", *response.Id)
}

func deleteReservation(ctx context.Context, client core.ComputeClient, reservationID string) {
	fmt.Printf("Deleting capacity reservation: %s\n", reservationID)

	request := core.DeleteComputeCapacityReservationRequest{
		CapacityReservationId: common.String(reservationID),
	}

	_, err := client.DeleteComputeCapacityReservation(ctx, request)
	if err != nil {
		log.Fatalf("Failed to delete capacity reservation: %v", err)
	}

	fmt.Println("✓ Capacity reservation deleted successfully!")
}

func isFlexibleShape(shape string) bool {
	flexibleShapes := []string{
		"VM.Standard.E3.Flex",
		"VM.Standard.E4.Flex",
		"VM.Standard.E5.Flex",
		"VM.Standard.A1.Flex",
		"VM.Optimized3.Flex",
	}
	
	for _, flex := range flexibleShapes {
		if shape == flex {
			return true
		}
	}
	return false
}
