package main

import (
	"context"
	"flag"
	"fmt"
	"log"
)

func main() {
	// Command-line flags
	var (
		configFile      = flag.String("config", "config.yaml", "Path to configuration file")
		instanceCount   = flag.Int("count", 0, "Number of instances in the pool (overrides config)")
		compartmentID   = flag.String("compartment", "", "Compartment OCID (overrides config)")
		displayName     = flag.String("name", "", "Instance pool display name (overrides config)")
		instancePoolID  = flag.String("pool-id", "", "Instance pool ID to scale (optional)")
		action          = flag.String("action", "create", "Action to perform: create, scale, terminate")
	)
	flag.Parse()

	// Load configuration
	config, err := LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Override config with command-line flags if provided
	if *instanceCount > 0 {
		config.InstancePool.Size = *instanceCount
	}
	if *compartmentID != "" {
		config.CompartmentID = *compartmentID
	}
	if *displayName != "" {
		config.InstancePool.DisplayName = *displayName
	}

	// Initialize OCI client
	client, err := NewOCIClient(config)
	if err != nil {
		log.Fatalf("Failed to initialize OCI client: %v", err)
	}

	ctx := context.Background()

	// Perform action
	switch *action {
	case "create":
		fmt.Printf("Creating instance pool with %d instances...\n", config.InstancePool.Size)
		pool, err := client.CreateInstancePool(ctx, config)
		if err != nil {
			log.Fatalf("Failed to create instance pool: %v", err)
		}
		fmt.Printf("Successfully created instance pool: %s (ID: %s)\n", *pool.DisplayName, *pool.Id)
		fmt.Printf("Instance pool is now provisioning. Check OCI console for status.\n")

	case "scale":
		if *instancePoolID == "" {
			log.Fatal("--pool-id is required for scale action")
		}
		fmt.Printf("Scaling instance pool %s to %d instances...\n", *instancePoolID, config.InstancePool.Size)
		err := client.ScaleInstancePool(ctx, *instancePoolID, config.InstancePool.Size)
		if err != nil {
			log.Fatalf("Failed to scale instance pool: %v", err)
		}
		fmt.Printf("Successfully scaled instance pool to %d instances\n", config.InstancePool.Size)

	case "terminate":
		if *instancePoolID == "" {
			log.Fatal("--pool-id is required for terminate action")
		}
		fmt.Printf("Terminating instance pool %s...\n", *instancePoolID)
		err := client.TerminateInstancePool(ctx, *instancePoolID)
		if err != nil {
			log.Fatalf("Failed to terminate instance pool: %v", err)
		}
		fmt.Printf("Successfully terminated instance pool\n")

	default:
		log.Fatalf("Unknown action: %s. Valid actions: create, scale, terminate", *action)
	}
}
