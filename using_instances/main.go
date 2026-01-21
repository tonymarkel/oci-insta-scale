package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
)

func main() {
	var (
		numInstances      = flag.Int("instances", 1, "Number of instances to create")
		displayName       = flag.String("name", "oci-instance", "Base name for instances")
		imageID           = flag.String("image", "", "Image ID (required)")
		shape             = flag.String("shape", "VM.Standard.E4.Flex", "Instance shape")
		subnetID          = flag.String("subnet", "", "Subnet ID (required)")
		compartmentID     = flag.String("compartment", "", "Compartment ID (required)")
		availabilityDomain = flag.String("ad", "", "Availability Domain (required)")
		outputFile        = flag.String("output", "instances.txt", "Output file for instance OCIDs")
	)
	flag.Parse()

	if *imageID == "" || *subnetID == "" || *compartmentID == "" || *availabilityDomain == "" {
		fmt.Println("Error: image, subnet, compartment, and ad flags are required")
		flag.PrintDefaults()
		return
	}

	// Create OCI client
	ctx := context.Background()
	configProvider := common.DefaultConfigProvider()
	client, err := core.NewComputeClientWithConfigurationProvider(configProvider)
	if err != nil {
		fmt.Printf("Error creating compute client: %v\n", err)
		return
	}

	fmt.Printf("Creating %d instance(s) in parallel...\n", *numInstances)

	// Create instances in parallel
	results := make(chan InstanceResult, *numInstances)
	var wg sync.WaitGroup

	for i := 1; i <= *numInstances; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			result := createInstance(ctx, client, InstanceConfig{
				CompartmentID:      *compartmentID,
				DisplayName:        fmt.Sprintf("%s-%d", *displayName, index),
				ImageID:            *imageID,
				Shape:              *shape,
				SubnetID:           *subnetID,
				AvailabilityDomain: *availabilityDomain,
			})
			results <- result
		}(i)
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and display results
	var instanceIDs []string
	successCount := 0
	for result := range results {
		if result.Error != nil {
			fmt.Printf("❌ Failed to create %s: %v\n", result.InstanceName, result.Error)
		} else {
			if result.RunningAt != nil {
				dur := result.RunningAt.Sub(result.LaunchStarted).Round(time.Second)
				fmt.Printf("✓ Successfully created %s (ID: %s) | launch: %s | running: %s | ready in: %s\n",
					result.InstanceName,
					result.InstanceID,
					result.LaunchStarted.Format(time.RFC3339),
					result.RunningAt.Format(time.RFC3339),
					dur,
				)
			} else {
				fmt.Printf("✓ Successfully created %s (ID: %s) | launch: %s\n",
					result.InstanceName,
					result.InstanceID,
					result.LaunchStarted.Format(time.RFC3339),
				)
			}
			instanceIDs = append(instanceIDs, result.InstanceID)
			successCount++
		}
	}

	fmt.Printf("\nSummary: %d/%d instances created successfully\n", successCount, *numInstances)

	// Write instance IDs to file
	if successCount > 0 {
		err := writeInstancesToFile(*outputFile, instanceIDs)
		if err != nil {
			fmt.Printf("Error writing instances to file: %v\n", err)
		} else {
			fmt.Printf("Instance OCIDs written to %s\n", *outputFile)
		}
	}
}

type InstanceConfig struct {
	CompartmentID      string
	DisplayName        string
	ImageID            string
	Shape              string
	SubnetID           string
	AvailabilityDomain string
}

type InstanceResult struct {
	InstanceName string
	InstanceID   string
	Error        error
	LaunchStarted time.Time
	RunningAt    *time.Time
}

func createInstance(ctx context.Context, client core.ComputeClient, config InstanceConfig) InstanceResult {
	launchStarted := time.Now().UTC()

	// Create launch instance details
	launchDetails := core.LaunchInstanceDetails{
		CompartmentId:      common.String(config.CompartmentID),
		DisplayName:        common.String(config.DisplayName),
		ImageId:            common.String(config.ImageID),
		Shape:              common.String(config.Shape),
		AvailabilityDomain: common.String(config.AvailabilityDomain),
		ShapeConfig: &core.LaunchInstanceShapeConfigDetails{
			Ocpus:       common.Float32(1.0),
			MemoryInGBs: common.Float32(8.0),
		},
		CreateVnicDetails: &core.CreateVnicDetails{
			SubnetId:            common.String(config.SubnetID),
			AssignPublicIp:      common.Bool(true),
			SkipSourceDestCheck: common.Bool(false),
		},
	}

	// Launch the instance
	request := core.LaunchInstanceRequest{
		LaunchInstanceDetails: launchDetails,
	}

	response, err := client.LaunchInstance(ctx, request)
	if err != nil {
		return InstanceResult{
			InstanceName: config.DisplayName,
			Error:        fmt.Errorf("launch failed: %w", err),
		}
	}

	result := InstanceResult{
		InstanceName:  config.DisplayName,
		InstanceID:    *response.Id,
		LaunchStarted: launchStarted,
	}

	runningAt, err := waitForInstanceRunning(ctx, client, result.InstanceID, 10*time.Second, 30*time.Minute)
	if err != nil {
		result.Error = fmt.Errorf("launch wait failed: %w", err)
		return result
	}

	result.RunningAt = &runningAt
	return result
}

func waitForInstanceRunning(ctx context.Context, client core.ComputeClient, instanceID string, interval time.Duration, maxWait time.Duration) (time.Time, error) {
	ctxWait, cancel := context.WithTimeout(ctx, maxWait)
	defer cancel()

	// Immediate check before waiting
	if t, done, err := checkInstanceRunning(ctxWait, client, instanceID); done || err != nil {
		return t, err
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctxWait.Done():
			return time.Time{}, fmt.Errorf("timeout waiting for running state: %w", ctxWait.Err())
		case <-ticker.C:
			t, done, err := checkInstanceRunning(ctxWait, client, instanceID)
			if err != nil {
				return time.Time{}, err
			}
			if done {
				return t, nil
			}
		}
	}
}

func checkInstanceRunning(ctx context.Context, client core.ComputeClient, instanceID string) (time.Time, bool, error) {
	resp, err := client.GetInstance(ctx, core.GetInstanceRequest{InstanceId: common.String(instanceID)})
	if err != nil {
		return time.Time{}, false, fmt.Errorf("get instance failed: %w", err)
	}

	switch resp.Instance.LifecycleState {
	case core.InstanceLifecycleStateRunning:
		return time.Now().UTC(), true, nil
	case core.InstanceLifecycleStateTerminated, core.InstanceLifecycleStateTerminating, core.InstanceLifecycleStateStopped:
		return time.Time{}, false, fmt.Errorf("instance entered terminal state: %s", resp.Instance.LifecycleState)
	default:
		return time.Time{}, false, nil
	}
}

func writeInstancesToFile(filename string, instanceIDs []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, id := range instanceIDs {
		_, err := writer.WriteString(id + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
