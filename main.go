package main

import (
	"bufio"
	"fmt"
	"net/http"
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

// Result holds the structural performance data collected from a target URL
type Result struct {
	URL        string
	StatusCode int
	Latency    time.Duration
	Err        error
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

// pingURL performs a fast HTTP GET request, explicitly tracking precise execution duration
func pingURL(url string, client *http.Client) Result {
	start := time.Now()
	resp, err := client.Get(url)
	latency := time.Since(start)

	if err != nil {
		return Result{URL: url, StatusCode: 0, Latency: latency, Err: err}
	}
	defer resp.Body.Close()

	return Result{URL: url, StatusCode: resp.StatusCode, Latency: latency, Err: nil}
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

	fmt.Printf("%-35s %-12s %-12s\n", "TARGET ENDPOINT", "STATUS", "LATENCY")
	fmt.Println(strings.Repeat("-", 65))

	for res := range resultChan {
		if res.Err != nil {
			fmt.Printf("%-35s %s%-12s%s %-12s\n", res.URL, Red, "FAILED", Reset, res.Latency.Round(time.Millisecond))
			failedCount++
		} else {
			statusStr := fmt.Sprintf("HTTP %d", res.StatusCode)
			if res.StatusCode == 200 {
				fmt.Printf("%-35s %s%-12s%s %-12s\n", res.URL, Green, statusStr, Reset, res.Latency.Round(time.Millisecond))
				successCount++
			} else {
				fmt.Printf("%-35s %s%-12s%s %-12s\n", res.URL, Yellow, statusStr, Reset, res.Latency.Round(time.Millisecond))
				successCount++
			}
		}
	}

	totalDuration := time.Since(startTime)
	fmt.Println(strings.Repeat("-", 65))
	fmt.Printf("\n📊 Execution Summary:\n")
	fmt.Printf("• Total Targets Processed: %d\n", len(urls))
	fmt.Printf("• Successful Enpoints:    %s%d%s\n", Green, successCount, Reset)
	fmt.Printf("• Failed/Unreachable:     %s%d%s\n", Red, failedCount, Reset)
	fmt.Printf("• Total Engine Wall Time: %s%v%s\n", Cyan, totalDuration.Round(time.Millisecond), Reset)
}

