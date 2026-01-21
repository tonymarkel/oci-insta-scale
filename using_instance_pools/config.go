package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	// OCI Authentication
	TenancyOCID     string `yaml:"tenancy_ocid"`
	UserOCID        string `yaml:"user_ocid"`
	Fingerprint     string `yaml:"fingerprint"`
	PrivateKeyPath  string `yaml:"private_key_path"`
	Region          string `yaml:"region"`
	
	// Instance Pool Configuration
	CompartmentID   string           `yaml:"compartment_id"`
	InstancePool    InstancePoolConfig `yaml:"instance_pool"`
}

// InstancePoolConfig defines the instance pool settings
type InstancePoolConfig struct {
	DisplayName           string                     `yaml:"display_name"`
	Size                  int                        `yaml:"size"`
	InstanceConfiguration InstanceConfigurationSpec  `yaml:"instance_configuration"`
	Placement             []PlacementConfig          `yaml:"placement"`
	LoadBalancers         []LoadBalancerConfig       `yaml:"load_balancers,omitempty"`
}

// InstanceConfigurationSpec defines the VM configuration
type InstanceConfigurationSpec struct {
	DisplayName       string            `yaml:"display_name"`
	Shape             string            `yaml:"shape"`
	ShapeConfig       ShapeConfig       `yaml:"shape_config,omitempty"`
	ImageID           string            `yaml:"image_id"`
	SubnetID          string            `yaml:"subnet_id"`
	AssignPublicIP    bool              `yaml:"assign_public_ip"`
	SSHAuthorizedKeys string            `yaml:"ssh_authorized_keys,omitempty"`
	UserData          string            `yaml:"user_data,omitempty"`
	Metadata          map[string]string `yaml:"metadata,omitempty"`
	FreeformTags      map[string]string `yaml:"freeform_tags,omitempty"`
	DefinedTags       map[string]map[string]interface{} `yaml:"defined_tags,omitempty"`
}

// ShapeConfig defines flexible shape configuration (for flex shapes)
type ShapeConfig struct {
	Ocpus       float32 `yaml:"ocpus,omitempty"`
	MemoryInGBs float32 `yaml:"memory_in_gbs,omitempty"`
}

// PlacementConfig defines availability domain and fault domain placement
type PlacementConfig struct {
	AvailabilityDomain string `yaml:"availability_domain"`
	FaultDomains       []string `yaml:"fault_domains,omitempty"`
}

// LoadBalancerConfig defines load balancer attachment
type LoadBalancerConfig struct {
	LoadBalancerID string `yaml:"load_balancer_id"`
	BackendSetName string `yaml:"backend_set_name"`
	Port           int    `yaml:"port"`
	VnicSelection  string `yaml:"vnic_selection"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate required fields
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// Validate checks if all required configuration fields are present
func (c *Config) Validate() error {
	if c.TenancyOCID == "" {
		return fmt.Errorf("tenancy_ocid is required")
	}
	if c.UserOCID == "" {
		return fmt.Errorf("user_ocid is required")
	}
	if c.Fingerprint == "" {
		return fmt.Errorf("fingerprint is required")
	}
	if c.PrivateKeyPath == "" {
		return fmt.Errorf("private_key_path is required")
	}
	if c.Region == "" {
		return fmt.Errorf("region is required")
	}
	if c.CompartmentID == "" {
		return fmt.Errorf("compartment_id is required")
	}
	// Note: Size can be 0 in config if overridden by command-line flag
	if c.InstancePool.InstanceConfiguration.Shape == "" {
		return fmt.Errorf("instance_pool.instance_configuration.shape is required")
	}
	if c.InstancePool.InstanceConfiguration.ImageID == "" {
		return fmt.Errorf("instance_pool.instance_configuration.image_id is required")
	}
	if c.InstancePool.InstanceConfiguration.SubnetID == "" {
		return fmt.Errorf("instance_pool.instance_configuration.subnet_id is required")
	}
	if len(c.InstancePool.Placement) == 0 {
		return fmt.Errorf("at least one placement configuration is required")
	}

	return nil
}
