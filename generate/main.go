package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

func main() {
	// Define flags
	outputDir := flag.String("output", "data/services", "Output directory for YAML files")
	jsonInputFile := flag.String("json-input-file", "", "Path to save the downloaded JSON file")
	forceFlag := flag.Bool("force", false, "Force generation even if no changes detected")
	verboseFlag := flag.Bool("verbose", false, "Enable verbose logging")

	flag.Parse()

	// Configure logging
	if *verboseFlag {
		log.SetLevel(log.DebugLevel)
	}

	log.Info("Azure IP Ranges Generator")

	if *jsonInputFile == "" {
		log.Fatal("json-input-file flag is required")
	}

	var jsonData []byte
	var err error

	// Download the JSON
	downloader := NewDownloader()
	url, err := downloader.FindJSONURL()
	if err != nil {
		log.Fatal("Failed to find JSON URL", "error", err)
	}

	jsonData, err = downloader.DownloadJSON(url)
	if err != nil {
		log.Fatal("Failed to download JSON", "error", err)
	}

	// Check for changes
	detector := NewChangeDetector(*jsonInputFile)
	hasChanges, err := detector.HasChanges(jsonData)
	if err != nil {
		log.Fatal("Failed to check for changes", "error", err)
	}

	if !hasChanges && !*forceFlag {
		log.Info("No changes detected, skipping generation")
		return
	}

	// Get change details
	details, err := detector.GetChangeDetails(jsonData)
	if err != nil {
		log.Warn("Failed to get change details", "error", err)
	} else {
		logChangeDetails(details)
	}

	// Save the JSON
	if err := os.WriteFile(*jsonInputFile, jsonData, 0644); err != nil {
		log.Fatal("Failed to save JSON file", "error", err)
	}
	log.Info("Saved JSON file", "path", *jsonInputFile)

	// Generate YAML files
	generator := NewGenerator(*outputDir)
	if err := generator.Generate(jsonData); err != nil {
		log.Fatal("Failed to generate YAML files", "error", err)
	}

	log.Info("Process complete")
}

func logChangeDetails(details *ChangeDetails) {
	if details.IsNew {
		log.Info("This is a new dataset")
		return
	}

	log.Info("Change details",
		"oldChangeNumber", details.OldChangeNumber,
		"newChangeNumber", details.NewChangeNumber,
		"oldServiceCount", details.ServiceCountOld,
		"newServiceCount", details.ServiceCountNew)

	if len(details.ServicesAdded) > 0 {
		log.Info("Services added", "count", len(details.ServicesAdded), "services", fmt.Sprintf("%v", details.ServicesAdded))
	}
	if len(details.ServicesRemoved) > 0 {
		log.Info("Services removed", "count", len(details.ServicesRemoved), "services", fmt.Sprintf("%v", details.ServicesRemoved))
	}
	if len(details.ServicesModified) > 0 {
		log.Info("Services modified", "count", len(details.ServicesModified))
		if len(details.ServicesModified) <= 10 {
			log.Debug("Modified services", "services", fmt.Sprintf("%v", details.ServicesModified))
		}
	}
}
