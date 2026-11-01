package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// ANSI Color Escape Codes for clean terminal output
const (
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
)

// Result holds the enriched performance data collected from a target URL
type Result struct {
	URL           string
	StatusCode    int
	Latency       time.Duration
	ServerHeader  string
	ContentLength int64
	ResolvedIP    string
	Err           error
}

// readURLs opens the input text file and extracts non-empty target strings
func readURLs(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			urls = append(urls, line)
		}
	}
	return urls, scanner.Err()
}

// pingURL performs a fast HTTP GET request and resolves the network endpoint IP address
func pingURL(targetURL string, client *http.Client) Result {
	start := time.Now()

	// Resolve IP address from the URL hostname
	ipAddress := "Unknown"
	if parsedURL, err := url.Parse(targetURL); err == nil {
		if ips, err := net.LookupIP(parsedURL.Hostname()); err == nil && len(ips) > 0 {
			ipAddress = ips[0].String()
		}
	}

	resp, err := client.Get(targetURL)
	latency := time.Since(start)

	if err != nil {
		return Result{URL: targetURL, StatusCode: 0, Latency: latency, ResolvedIP: ipAddress, Err: err}
	}
	defer resp.Body.Close()

	// Extract standard server infrastructure headers
	server := resp.Header.Get("Server")
	if server == "" {
		server = "N/A"
	}

	return Result{
		URL:           targetURL,
		StatusCode:    resp.StatusCode,
		Latency:       latency,
		ServerHeader:  server,
		ContentLength: resp.ContentLength,
		ResolvedIP:    ipAddress,
		Err:           nil,
	}
}

// MaxConcurrency limits the number of simultaneous network workers running at one time
const MaxConcurrency = 3

func main() {
	targetFile := "urls.txt"
	fmt.Printf("%s🔍 Initializing Go Concurrency Engine... Reading '%s'%s\n\n", Cyan, targetFile, Reset)

	urls, err := readURLs(targetFile)
	if err != nil {
		fmt.Printf("%s[FATAL]%s Failed to read target file: %v\n", Red, Reset, err)
		os.Exit(1)
	}

	if len(urls) == 0 {
		fmt.Println("⚠ No target URLs found to process.")
		return
	}

	// Create a thread-safe HTTP Client with a strict global timeout threshold
	client := &http.Client{
		Timeout: 6 * time.Second,
	}

	// Initialize synchronization primitives
	var wg sync.WaitGroup
	resultChan := make(chan Result, len(urls))

	// Semaphore channel to enforce bounded concurrency limits
	sem := make(chan struct{}, MaxConcurrency)

	startTime := time.Now()

	// Spin up a separate Goroutine for each URL concurrently
	for _, url := range urls {
		wg.Add(1)
		go func(targetURL string) {
			// Acquire token: blocks execution if the 'sem' channel buffer is fully saturated
			sem <- struct{}{}

			// Release token when worker completes, freeing up space for subsequent threads
			defer func() { <-sem }()
			defer wg.Done()

			// Execute the network probe and stream the structural data into our channel
			resultChan <- pingURL(targetURL, client)
		}(url)
	}

	// Monitor the worker threads and safely close the channel once execution finishes
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Read and process telemetry live as it flows down the channel
	successCount := 0
	failedCount := 0

	// Enriched table headers matching our new data streams
	fmt.Printf("%-24s %-10s %-10s %-16s %-12s %s\n", "TARGET ENDPOINT", "STATUS", "LATENCY", "RESOLVED IP", "SIZE(BYTES)", "SERVER")
	fmt.Println(strings.Repeat("-", 90))

	for res := range resultChan {
		// Truncate long URLs gracefully for table layout alignment
		displayURL := res.URL
		if len(displayURL) > 24 {
			displayURL = displayURL[:21] + "..."
		}

		if res.Err != nil {
			fmt.Printf("%-24s %s%-10s%s %-10s %-16s %-12s %s\n",
				   displayURL, Red, "FAILED", Reset, res.Latency.Round(time.Millisecond), res.ResolvedIP, "0", "N/A")
			failedCount++
		} else {
			statusStr := fmt.Sprintf("HTTP %d", res.StatusCode)
			sizeStr := fmt.Sprintf("%d", res.ContentLength)
			if res.ContentLength < 0 {
				sizeStr = "Chunked" // Handle dynamic streaming transfer-encodings
			}

			colorCode := Yellow
			if res.StatusCode == 200 {
				colorCode = Green
				successCount++
			} else {
				successCount++
			}

			fmt.Printf("%-24s %s%-10s%s %-10s %-16s %-12s %s\n",
				   displayURL, colorCode, statusStr, Reset, res.Latency.Round(time.Millisecond), res.ResolvedIP, sizeStr, res.ServerHeader)
		}
	}

	totalDuration := time.Since(startTime)
	fmt.Println(strings.Repeat("-", 90))
	fmt.Printf("\n📊 Execution Summary:\n")
	fmt.Printf("• Total Targets Processed: %d\n", len(urls))
	fmt.Printf("• Successful Endpoints:   %s%d%s\n", Green, successCount, Reset)
	fmt.Printf("• Failed/Unreachable:     %s%d%s\n", Red, failedCount, Reset)
	fmt.Printf("• Total Engine Wall Time: %s%v%s\n", Cyan, totalDuration.Round(time.Millisecond), Reset)
}
