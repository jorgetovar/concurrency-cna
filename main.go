package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func responseSize(url string) int {
	fmt.Println("Getting", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return 0
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return 0
	}
	return len(body)
}

func responseSizeWithChannel(url string, channel chan Page) {
	fmt.Println("Getting", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	channel <- Page{URL: url, Size: len(body)}
}

type Page struct {
	URL  string
	Size int
}

func main() {

	start := time.Now()
	fmt.Println("Hello AWS Community!")
	const responseFormat = "Response size: %d\n"
	fmt.Printf(responseFormat, responseSize("https://example.com"))
	fmt.Printf(responseFormat, responseSize("https://google.com"))
	fmt.Printf(responseFormat, responseSize("https://github.com/jorgetovar"))
	fmt.Printf("Time taken 3 URLs %v (Sync)\n", time.Since(start))
	startGoroutines := time.Now()
	urls := []string{"https://example.com", "https://google.com", "https://github.com/jorgetovar"}
	pages := make(chan Page)
	for _, url := range urls {
		go responseSizeWithChannel(url, pages)
	}

	for i := 0; i < len(urls); i++ {
		page := <-pages
		fmt.Printf("Response size: %d for URL %s \n", page.Size, page.URL)
	}
	fmt.Printf("Time taken 3 URLs %v (Goroutines & Channels)\n ", time.Since(startGoroutines))

}
