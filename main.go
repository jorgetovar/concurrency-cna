package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Page represents the result of fetching a URL
type Page struct {
	URL      string
	Size     int
	Duration time.Duration
	Error    error
}

// fetchURL handles the HTTP request and returns a Page
func fetchURL(url string) Page {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return Page{URL: url, Error: err}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Page{URL: url, Error: err}
	}

	return Page{
		URL:      url,
		Size:     len(body),
		Duration: time.Since(start),
	}
}

// fetchURLsSequential demonstrates fetching URLs without goroutines
func fetchURLsSequential(urls []string) []Page {
	results := make([]Page, len(urls))
	for i, url := range urls {
		results[i] = fetchURL(url)
	}
	return results
}

// fetchURLsConcurrent demonstrates basic channel usage with goroutines
func fetchURLsConcurrent(urls []string) []Page {
	results := make(chan Page, len(urls))
	var wg sync.WaitGroup

	// Launch a goroutine for each URL
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			results <- fetchURL(url)
		}(url)
	}

	// Close results channel when all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var pages []Page
	for page := range results {
		pages = append(pages, page)
	}
	return pages
}

// workerPool demonstrates a worker pool pattern with channels
func workerPool(urls []string, numWorkers int) []Page {
	jobs := make(chan string, len(urls))  // Channel for URLs to process
	results := make(chan Page, len(urls)) // Channel for results
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range jobs {
				results <- fetchURL(url)
			}
		}()
	}

	// Send jobs
	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	// Wait for workers to finish and close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var pages []Page
	for page := range results {
		pages = append(pages, page)
	}
	return pages
}

func main() {
	urls := []string{
		"https://example.com",
		"https://google.com",
		"https://github.com",
		"https://go.dev",
		"https://aws.amazon.com",
	}

	// Example 1: Sequential execution
	fmt.Println("Sequential Fetching:")
	start := time.Now()
	sequential := fetchURLsSequential(urls)
	fmt.Printf("\nSequential execution took: %v\n", time.Since(start))
	const formatMessage = "Error fetching %s: %v\n"
	const bytesMessage = "%s: %d bytes in %v\n"
	for _, page := range sequential {
		if page.Error != nil {
			fmt.Printf(formatMessage, page.URL, page.Error)
		} else {
			fmt.Printf(bytesMessage, page.URL, page.Size, page.Duration)
		}
	}

	// Example 2: Concurrent execution with goroutines
	fmt.Println("\nConcurrent Fetching (unlimited goroutines):")
	start = time.Now()
	concurrent := fetchURLsConcurrent(urls)
	fmt.Printf("\nConcurrent execution took: %v\n", time.Since(start))
	for _, page := range concurrent {
		if page.Error != nil {
			fmt.Printf(formatMessage, page.URL, page.Error)
		} else {
			fmt.Printf(bytesMessage, page.URL, page.Size, page.Duration)
		}
	}

	// Example 3: Worker pool for controlled concurrency
	fmt.Println("\nWorker Pool (3 workers):")
	start = time.Now()
	workerResults := workerPool(urls, 3)
	fmt.Printf("\nWorker pool execution took: %v\n", time.Since(start))
	for _, page := range workerResults {
		if page.Error != nil {
			fmt.Printf(formatMessage, page.URL, page.Error)
		} else {
			fmt.Printf(bytesMessage, page.URL, page.Size, page.Duration)
		}
	}
}
