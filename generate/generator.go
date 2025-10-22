package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
)

// ServiceTags represents the top-level JSON structure
type ServiceTags struct {
	ChangeNumber int       `json:"changeNumber"`
	Cloud        string    `json:"cloud"`
	Values       []Service `json:"values"`
}

// Service represents a single Azure service
type Service struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Properties ServiceProperties `json:"properties"`
}

// ServiceProperties contains service metadata and IP ranges
type ServiceProperties struct {
	ChangeNumber    int      `json:"changeNumber"`
	Region          string   `json:"region"`
	Platform        string   `json:"platform"`
	SystemService   string   `json:"systemService"`
	AddressPrefixes []string `json:"addressPrefixes"`
	NetworkFeatures []string `json:"networkFeatures"`
}

// ServiceYAML represents the YAML output structure
type ServiceYAML struct {
	ID              string              `yaml:"id"`
	Name            string              `yaml:"name"`
	Metadata        MetadataYAML        `yaml:"metadata"`
	AddressPrefixes AddressPrefixesYAML `yaml:"address_prefixes"`
}

// MetadataYAML contains service metadata
type MetadataYAML struct {
	ChangeNumber       int      `yaml:"change_number"`
	Region             string   `yaml:"region"`
	Platform           string   `yaml:"platform"`
	SystemService      string   `yaml:"system_service"`
	NetworkFeatures    []string `yaml:"network_features"`
	GlobalChangeNumber int      `yaml:"global_change_number"`
	Cloud              string   `yaml:"cloud"`
}

// AddressPrefixesYAML contains categorized IP prefixes
type AddressPrefixesYAML struct {
	All    []string   `yaml:"all"`
	IPv4   []string   `yaml:"ipv4"`
	IPv6   []string   `yaml:"ipv6"`
	Counts CountsYAML `yaml:"counts"`
}

// CountsYAML contains counts of IP prefixes
type CountsYAML struct {
	Total int `yaml:"total"`
	IPv4  int `yaml:"ipv4"`
	IPv6  int `yaml:"ipv6"`
}

// GlobalMetadataYAML represents the metadata.yaml file
type GlobalMetadataYAML struct {
	ChangeNumber int       `yaml:"change_number"`
	Cloud        string    `yaml:"cloud"`
	ServiceCount int       `yaml:"service_count"`
	GeneratedAt  time.Time `yaml:"generated_at"`
}

// Generator handles YAML generation from JSON data
type Generator struct {
	outputDir string
}

// NewGenerator creates a new generator
func NewGenerator(outputDir string) *Generator {
	return &Generator{
		outputDir: outputDir,
	}
}

// Generate processes the JSON data and generates YAML files
func (g *Generator) Generate(jsonData []byte) error {
	var serviceTags ServiceTags
	if err := json.Unmarshal(jsonData, &serviceTags); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Create output directory
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	log.Info("Processing services",
		"count", len(serviceTags.Values),
		"changeNumber", serviceTags.ChangeNumber,
		"cloud", serviceTags.Cloud)

	// Generate global metadata file
	metadata := GlobalMetadataYAML{
		ChangeNumber: serviceTags.ChangeNumber,
		Cloud:        serviceTags.Cloud,
		ServiceCount: len(serviceTags.Values),
		GeneratedAt:  time.Now().UTC(),
	}

	metadataFile := filepath.Join(g.outputDir, "metadata.yaml")
	if err := writeYAML(metadataFile, metadata); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}
	log.Info("Generated metadata file", "file", metadataFile)

	// Generate YAML file for each service
	successCount := 0
	for _, service := range serviceTags.Values {
		if service.ID == "" {
			continue
		}

		serviceYAML := generateServiceYAML(service, serviceTags.ChangeNumber, serviceTags.Cloud)

		filename := fmt.Sprintf("%s.yaml", sanitizeFilename(service.ID))
		outputFile := filepath.Join(g.outputDir, filename)

		if err := writeYAML(outputFile, serviceYAML); err != nil {
			log.Warn("Failed to write service file", "filename", filename, "error", err)
			continue
		}

		successCount++
		log.Debug("Generated service file",
			"filename", filename,
			"prefixes", serviceYAML.AddressPrefixes.Counts.Total,
			"ipv4", serviceYAML.AddressPrefixes.Counts.IPv4,
			"ipv6", serviceYAML.AddressPrefixes.Counts.IPv6)
	}

	log.Info("YAML generation complete",
		"totalServices", len(serviceTags.Values),
		"successfulFiles", successCount,
		"outputDir", g.outputDir)

	return nil
}

// isIPv4 checks if an IP prefix is IPv4
func isIPv4(prefix string) bool {
	return strings.Contains(prefix, ".") && !strings.Contains(prefix, ":")
}

// isIPv6 checks if an IP prefix is IPv6
func isIPv6(prefix string) bool {
	return strings.Contains(prefix, ":")
}

// categorizeIPPrefixes splits prefixes into IPv4 and IPv6
func categorizeIPPrefixes(prefixes []string) ([]string, []string) {
	var ipv4, ipv6 []string
	for _, prefix := range prefixes {
		if isIPv4(prefix) {
			ipv4 = append(ipv4, prefix)
		} else if isIPv6(prefix) {
			ipv6 = append(ipv6, prefix)
		}
	}
	return ipv4, ipv6
}

// sanitizeFilename converts service ID to a safe filename
func sanitizeFilename(serviceID string) string {
	return strings.ToLower(strings.ReplaceAll(serviceID, ".", "_"))
}

// generateServiceYAML creates the YAML structure for a service
func generateServiceYAML(service Service, globalChangeNumber int, cloud string) ServiceYAML {
	ipv4, ipv6 := categorizeIPPrefixes(service.Properties.AddressPrefixes)

	return ServiceYAML{
		ID:   service.ID,
		Name: service.Name,
		Metadata: MetadataYAML{
			ChangeNumber:       service.Properties.ChangeNumber,
			Region:             service.Properties.Region,
			Platform:           service.Properties.Platform,
			SystemService:      service.Properties.SystemService,
			NetworkFeatures:    service.Properties.NetworkFeatures,
			GlobalChangeNumber: globalChangeNumber,
			Cloud:              cloud,
		},
		AddressPrefixes: AddressPrefixesYAML{
			All:  service.Properties.AddressPrefixes,
			IPv4: ipv4,
			IPv6: ipv6,
			Counts: CountsYAML{
				Total: len(service.Properties.AddressPrefixes),
				IPv4:  len(ipv4),
				IPv6:  len(ipv6),
			},
		},
	}
}

// writeYAML writes data to a YAML file
func writeYAML(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	return encoder.Encode(data)
}
