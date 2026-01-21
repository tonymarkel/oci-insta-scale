package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
)

func init() {
	// Only run this if called with "terminate" as first argument
	if len(os.Args) > 1 && os.Args[1] == "terminate" {
		os.Args = append(os.Args[:1], os.Args[2:]...)
		runTerminate()
		os.Exit(0)
	}
}

func runTerminate() {
	var (
		inputFile   = flag.String("file", "instances.txt", "File containing instance OCIDs (one per line)")
		compartment = flag.String("compartment", "", "Compartment ID (required)")
		parallel    = flag.Int("parallel", 10, "Number of parallel termination operations")
	)
	flag.Parse()

	if *compartment == "" {
		fmt.Println("Error: compartment flag is required for termination")
		flag.PrintDefaults()
		return
	}

	// Read instance IDs from file
	instanceIDs, err := readInstancesFromFile(*inputFile)
	if err != nil {
		fmt.Printf("Error reading instances from file: %v\n", err)
		return
	}

	if len(instanceIDs) == 0 {
		fmt.Println("No instance IDs found in file")
		return
	}

	fmt.Printf("Found %d instances to terminate\n", len(instanceIDs))

	// Create OCI client
	ctx := context.Background()
	configProvider := common.DefaultConfigProvider()
	client, err := core.NewComputeClientWithConfigurationProvider(configProvider)
	if err != nil {
		fmt.Printf("Error creating compute client: %v\n", err)
		return
	}

	// Terminate instances with concurrency limit
	results := make(chan TerminationResult, len(instanceIDs))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, *parallel)

	for _, instanceID := range instanceIDs {
		time.Sleep(1500 * time.Millisecond) // Slight delay to avoid overwhelming the API
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			result := terminateInstance(ctx, client, id, *compartment)
			results <- result
		}(instanceID)
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and display results
	successCount := 0
	failureCount := 0
	for result := range results {
		if result.Error != nil {
			fmt.Printf("❌ Failed to terminate %s: %v\n", result.InstanceID, result.Error)
			failureCount++
		} else {
			fmt.Printf("✓ Successfully terminated %s\n", result.InstanceID)
			successCount++
		}
	}

	fmt.Printf("\nSummary: %d/%d instances terminated successfully\n", successCount, len(instanceIDs))
	if failureCount > 0 {
		fmt.Printf("Failures: %d\n", failureCount)
	}
}

type TerminationResult struct {
	InstanceID string
	Error      error
}

func terminateInstance(ctx context.Context, client core.ComputeClient, instanceID, compartmentID string) TerminationResult {
	request := core.TerminateInstanceRequest{
		InstanceId: common.String(instanceID),
	}

	_, err := client.TerminateInstance(ctx, request)
	if err != nil {
		return TerminationResult{
			InstanceID: instanceID,
			Error:      fmt.Errorf("terminate failed: %w", err),
		}
	}

	return TerminationResult{
		InstanceID: instanceID,
		Error:      nil,
	}
}

func readInstancesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var instanceIDs []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			instanceIDs = append(instanceIDs, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return instanceIDs, nil
}
