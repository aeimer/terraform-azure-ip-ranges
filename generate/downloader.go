package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/charmbracelet/log"
)

const (
	downloadPageURL = "https://www.microsoft.com/en-us/download/details.aspx?id=56519"
)

// Downloader handles downloading Azure IP ranges data
type Downloader struct {
	client *http.Client
}

// NewDownloader creates a new downloader
func NewDownloader() *Downloader {
	return &Downloader{
		client: &http.Client{},
	}
}

// FindJSONURL scrapes the download page to find the JSON file URL
func (d *Downloader) FindJSONURL() (string, error) {
	log.Info("Fetching download page", "url", downloadPageURL)

	resp, err := d.client.Get(downloadPageURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch download page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Find all links matching the pattern
	pattern := `href="(https://download\.microsoft\.com/download/[^"]+\.json)"`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(string(body), -1)

	if len(matches) == 0 {
		return "", fmt.Errorf("no JSON download link found on page")
	}

	if len(matches) > 1 {
		return "", fmt.Errorf("multiple JSON download links found (%d), expected exactly one", len(matches))
	}

	url := matches[0][1]
	// Decode HTML entities if needed
	url = strings.ReplaceAll(url, "&amp;", "&")

	log.Info("Found JSON download URL", "url", url)
	return url, nil
}

// DownloadJSON downloads the JSON file from the given URL
func (d *Downloader) DownloadJSON(url string) ([]byte, error) {
	log.Info("Downloading JSON file", "url", url)

	resp, err := d.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download JSON: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON data: %w", err)
	}

	log.Info("Successfully downloaded JSON", "size", len(data))
	return data, nil
}
