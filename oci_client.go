package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
)

// OCIClient wraps OCI SDK clients
type OCIClient struct {
	ComputeClient         core.ComputeClient
	ComputeManagementClient core.ComputeManagementClient
	Config                *Config
}

// NewOCIClient creates a new OCI client with authentication
func NewOCIClient(config *Config) (*OCIClient, error) {
	// Read the private key file
	privateKey, err := os.ReadFile(config.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	// Create configuration provider with the private key content
	configProvider := common.NewRawConfigurationProvider(
		config.TenancyOCID,
		config.UserOCID,
		config.Region,
		config.Fingerprint,
		string(privateKey),
		nil, // passphrase
	)

	// Create compute client
	computeClient, err := core.NewComputeClientWithConfigurationProvider(configProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute client: %w", err)
	}

	// Create compute management client
	computeMgmtClient, err := core.NewComputeManagementClientWithConfigurationProvider(configProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute management client: %w", err)
	}

	return &OCIClient{
		ComputeClient:         computeClient,
		ComputeManagementClient: computeMgmtClient,
		Config:                config,
	}, nil
}

// CreateInstancePool creates an instance pool with the specified configuration
func (c *OCIClient) CreateInstancePool(ctx context.Context, config *Config) (*core.InstancePool, error) {
	// Step 1: Create instance configuration
	fmt.Println("Creating instance configuration...")
	instanceConfig, err := c.createInstanceConfiguration(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance configuration: %w", err)
	}
	fmt.Printf("Instance configuration created: %s\n", *instanceConfig.Id)

	// Step 2: Build placement configurations
	placementConfigs := make([]core.CreateInstancePoolPlacementConfigurationDetails, 0, len(config.InstancePool.Placement))
	for _, placement := range config.InstancePool.Placement {
		placementConfig := core.CreateInstancePoolPlacementConfigurationDetails{
			AvailabilityDomain: common.String(placement.AvailabilityDomain),
			PrimarySubnetId:    common.String(config.InstancePool.InstanceConfiguration.SubnetID),
		}
		if len(placement.FaultDomains) > 0 {
			placementConfig.FaultDomains = placement.FaultDomains
		}
		placementConfigs = append(placementConfigs, placementConfig)
	}

	// Step 3: Build load balancer configurations (if any)
	var lbAttachments []core.AttachLoadBalancerDetails
	for _, lb := range config.InstancePool.LoadBalancers {
		lbAttachment := core.AttachLoadBalancerDetails{
			LoadBalancerId: common.String(lb.LoadBalancerID),
			BackendSetName: common.String(lb.BackendSetName),
			Port:           common.Int(lb.Port),
			VnicSelection:  common.String(lb.VnicSelection),
		}
		lbAttachments = append(lbAttachments, lbAttachment)
	}

	// Step 4: Create the instance pool
	fmt.Println("Creating instance pool...")
	displayName := config.InstancePool.DisplayName
	if displayName == "" {
		displayName = fmt.Sprintf("instance-pool-%d", time.Now().Unix())
	}

	createPoolReq := core.CreateInstancePoolRequest{
		CreateInstancePoolDetails: core.CreateInstancePoolDetails{
			CompartmentId:              common.String(config.CompartmentID),
			InstanceConfigurationId:    instanceConfig.Id,
			PlacementConfigurations:    placementConfigs,
			Size:                       common.Int(config.InstancePool.Size),
			DisplayName:                common.String(displayName),
			LoadBalancers:              lbAttachments,
		},
	}

	poolResp, err := c.ComputeManagementClient.CreateInstancePool(ctx, createPoolReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance pool: %w", err)
	}

	return &poolResp.InstancePool, nil
}

// createInstanceConfiguration creates an instance configuration from the config
func (c *OCIClient) createInstanceConfiguration(ctx context.Context, config *Config) (*core.InstanceConfiguration, error) {
	instConfig := config.InstancePool.InstanceConfiguration
	
	displayName := instConfig.DisplayName
	if displayName == "" {
		displayName = fmt.Sprintf("instance-config-%d", time.Now().Unix())
	}

	// Build instance details
	instanceDetails := core.ComputeInstanceDetails{
		LaunchDetails: &core.InstanceConfigurationLaunchInstanceDetails{
			CompartmentId: common.String(config.CompartmentID),
			Shape:         common.String(instConfig.Shape),
			CreateVnicDetails: &core.InstanceConfigurationCreateVnicDetails{
				SubnetId:       common.String(instConfig.SubnetID),
				AssignPublicIp: common.Bool(instConfig.AssignPublicIP),
			},
			SourceDetails: &core.InstanceConfigurationInstanceSourceViaImageDetails{
				ImageId: common.String(instConfig.ImageID),
			},
		},
	}

	// Add shape config for flexible shapes
	if instConfig.ShapeConfig.Ocpus > 0 || instConfig.ShapeConfig.MemoryInGBs > 0 {
		instanceDetails.LaunchDetails.ShapeConfig = &core.InstanceConfigurationLaunchInstanceShapeConfigDetails{}
		if instConfig.ShapeConfig.Ocpus > 0 {
			instanceDetails.LaunchDetails.ShapeConfig.Ocpus = common.Float32(instConfig.ShapeConfig.Ocpus)
		}
		if instConfig.ShapeConfig.MemoryInGBs > 0 {
			instanceDetails.LaunchDetails.ShapeConfig.MemoryInGBs = common.Float32(instConfig.ShapeConfig.MemoryInGBs)
		}
	}

	// Add SSH keys
	if instConfig.SSHAuthorizedKeys != "" {
		instanceDetails.LaunchDetails.Metadata = map[string]string{
			"ssh_authorized_keys": instConfig.SSHAuthorizedKeys,
		}
	}

	// Add user data
	if instConfig.UserData != "" {
		if instanceDetails.LaunchDetails.Metadata == nil {
			instanceDetails.LaunchDetails.Metadata = make(map[string]string)
		}
		instanceDetails.LaunchDetails.Metadata["user_data"] = instConfig.UserData
	}

	// Add custom metadata
	if len(instConfig.Metadata) > 0 {
		if instanceDetails.LaunchDetails.Metadata == nil {
			instanceDetails.LaunchDetails.Metadata = make(map[string]string)
		}
		for k, v := range instConfig.Metadata {
			instanceDetails.LaunchDetails.Metadata[k] = v
		}
	}

	// Add tags
	if len(instConfig.FreeformTags) > 0 {
		instanceDetails.LaunchDetails.FreeformTags = instConfig.FreeformTags
	}
	if len(instConfig.DefinedTags) > 0 {
		instanceDetails.LaunchDetails.DefinedTags = instConfig.DefinedTags
	}

	// Create the instance configuration
	createConfigReq := core.CreateInstanceConfigurationRequest{
		CreateInstanceConfiguration: core.CreateInstanceConfigurationDetails{
			CompartmentId:   common.String(config.CompartmentID),
			DisplayName:     common.String(displayName),
			InstanceDetails: instanceDetails,
		},
	}

	configResp, err := c.ComputeManagementClient.CreateInstanceConfiguration(ctx, createConfigReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance configuration: %w", err)
	}

	return &configResp.InstanceConfiguration, nil
}

// ScaleInstancePool scales an existing instance pool to a new size
func (c *OCIClient) ScaleInstancePool(ctx context.Context, instancePoolID string, newSize int) error {
	updateReq := core.UpdateInstancePoolRequest{
		InstancePoolId: common.String(instancePoolID),
		UpdateInstancePoolDetails: core.UpdateInstancePoolDetails{
			Size: common.Int(newSize),
		},
	}

	_, err := c.ComputeManagementClient.UpdateInstancePool(ctx, updateReq)
	if err != nil {
		return fmt.Errorf("failed to update instance pool: %w", err)
	}

	return nil
}

// TerminateInstancePool terminates an instance pool and all its instances
func (c *OCIClient) TerminateInstancePool(ctx context.Context, instancePoolID string) error {
	terminateReq := core.TerminateInstancePoolRequest{
		InstancePoolId: common.String(instancePoolID),
	}

	_, err := c.ComputeManagementClient.TerminateInstancePool(ctx, terminateReq)
	if err != nil {
		return fmt.Errorf("failed to terminate instance pool: %w", err)
	}

	return nil
}

// GetInstancePool retrieves details about an instance pool
func (c *OCIClient) GetInstancePool(ctx context.Context, instancePoolID string) (*core.InstancePool, error) {
	getReq := core.GetInstancePoolRequest{
		InstancePoolId: common.String(instancePoolID),
	}

	resp, err := c.ComputeManagementClient.GetInstancePool(ctx, getReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance pool: %w", err)
	}

	return &resp.InstancePool, nil
}

// ListInstancePoolInstances lists all instances in an instance pool
func (c *OCIClient) ListInstancePoolInstances(ctx context.Context, compartmentID, instancePoolID string) ([]core.InstanceSummary, error) {
	listReq := core.ListInstancePoolInstancesRequest{
		CompartmentId:  common.String(compartmentID),
		InstancePoolId: common.String(instancePoolID),
	}

	resp, err := c.ComputeManagementClient.ListInstancePoolInstances(ctx, listReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list instance pool instances: %w", err)
	}

	return resp.Items, nil
}
