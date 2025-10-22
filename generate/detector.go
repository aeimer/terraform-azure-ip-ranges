package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

// ChangeDetector detects changes between old and new JSON data
type ChangeDetector struct {
	oldDataPath string
}

// NewChangeDetector creates a new change detector
func NewChangeDetector(oldDataPath string) *ChangeDetector {
	return &ChangeDetector{
		oldDataPath: oldDataPath,
	}
}

// HasChanges checks if the new data differs from the old data
func (cd *ChangeDetector) HasChanges(newData []byte) (bool, error) {
	// Check if old file exists
	if _, err := os.Stat(cd.oldDataPath); os.IsNotExist(err) {
		log.Info("No previous data file found, treating as new data")
		return true, nil
	}

	oldData, err := os.ReadFile(cd.oldDataPath)
	if err != nil {
		return false, fmt.Errorf("failed to read old data: %w", err)
	}

	// Parse both files to compare change numbers
	var oldTags, newTags ServiceTags
	if err := json.Unmarshal(oldData, &oldTags); err != nil {
		log.Warn("Failed to parse old JSON, treating as changed", "error", err)
		return true, nil
	}
	if err := json.Unmarshal(newData, &newTags); err != nil {
		return false, fmt.Errorf("failed to parse new JSON: %w", err)
	}

	// Compare change numbers
	if oldTags.ChangeNumber != newTags.ChangeNumber {
		log.Info("Change number differs",
			"old", oldTags.ChangeNumber,
			"new", newTags.ChangeNumber)
		return true, nil
	}

	// If change numbers are the same, do a byte comparison as final check
	if !bytes.Equal(oldData, newData) {
		log.Info("Change numbers match but content differs, treating as changed")
		return true, nil
	}

	log.Info("No changes detected", "changeNumber", oldTags.ChangeNumber)
	return false, nil
}

// GetChangeDetails returns detailed information about what changed
func (cd *ChangeDetector) GetChangeDetails(newData []byte) (*ChangeDetails, error) {
	details := &ChangeDetails{
		IsNew: false,
	}

	// Check if old file exists
	if _, err := os.Stat(cd.oldDataPath); os.IsNotExist(err) {
		details.IsNew = true
		return details, nil
	}

	oldData, err := os.ReadFile(cd.oldDataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read old data: %w", err)
	}

	var oldTags, newTags ServiceTags
	if err := json.Unmarshal(oldData, &oldTags); err != nil {
		log.Warn("Failed to parse old JSON", "error", err)
		details.IsNew = true
		return details, nil
	}
	if err := json.Unmarshal(newData, &newTags); err != nil {
		return nil, fmt.Errorf("failed to parse new JSON: %w", err)
	}

	details.OldChangeNumber = oldTags.ChangeNumber
	details.NewChangeNumber = newTags.ChangeNumber
	details.ServiceCountOld = len(oldTags.Values)
	details.ServiceCountNew = len(newTags.Values)

	// Build maps for comparison
	oldServices := make(map[string]Service)
	for _, svc := range oldTags.Values {
		oldServices[svc.ID] = svc
	}

	newServices := make(map[string]Service)
	for _, svc := range newTags.Values {
		newServices[svc.ID] = svc
	}

	// Find added and removed services
	for id := range newServices {
		if _, exists := oldServices[id]; !exists {
			details.ServicesAdded = append(details.ServicesAdded, id)
		}
	}
	for id := range oldServices {
		if _, exists := newServices[id]; !exists {
			details.ServicesRemoved = append(details.ServicesRemoved, id)
		}
	}

	// Find modified services
	for id, newSvc := range newServices {
		if oldSvc, exists := oldServices[id]; exists {
			if !servicesEqual(oldSvc, newSvc) {
				details.ServicesModified = append(details.ServicesModified, id)
			}
		}
	}

	return details, nil
}

// ChangeDetails contains information about what changed
type ChangeDetails struct {
	IsNew            bool
	OldChangeNumber  int
	NewChangeNumber  int
	ServiceCountOld  int
	ServiceCountNew  int
	ServicesAdded    []string
	ServicesRemoved  []string
	ServicesModified []string
}

// servicesEqual compares two services for equality
func servicesEqual(a, b Service) bool {
	if a.ID != b.ID || a.Name != b.Name {
		return false
	}
	if len(a.Properties.AddressPrefixes) != len(b.Properties.AddressPrefixes) {
		return false
	}
	// For simplicity, we'll just compare change numbers
	return a.Properties.ChangeNumber == b.Properties.ChangeNumber
}
