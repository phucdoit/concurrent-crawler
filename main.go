package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
)

const (
	queueSize         = 10
	numOfCrawlWorkers = 5
	numOfWriteWorkers = 3
)

type CrawlData struct {
	Url  string
	Data string
	Err  error
}

func crawl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func crawlConsumer(wg *sync.WaitGroup, crawlCh <-chan string, writeCh chan<- CrawlData) {
	defer wg.Done()
	for v := range crawlCh {
		data, err := crawl(v)
		if err != nil {
			log.Println("Failed to crawl URL:", v, "Error:", err)
			writeCh <- CrawlData{Url: v, Err: err}
			continue
		}
		writeCh <- CrawlData{Url: v, Data: string(data)}
	}
}

func write(data CrawlData) error {
	if data.Err != nil {
		return data.Err
	}

	if err := os.MkdirAll("output", os.ModePerm); err != nil {
		return err
	}

	fileName := "output/" + normalizeFilename(data.Url)
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.WriteString(f, data.Data)
	return err
}

func writeConsumer(wg *sync.WaitGroup, writeCh <-chan CrawlData) {
	defer wg.Done()
	for v := range writeCh {
		if err := write(v); err != nil {
			log.Println("Failed to write URL:", v.Url, "Error:", err)
			continue
		}
	}
}

func CrawAndWrite(urls []string) {
	// Queue Crawl
	crawlCh := make(chan string, queueSize)
	// Queue Write
	writeCh := make(chan CrawlData, queueSize)

	// Start crawl workers
	var wgCrawl sync.WaitGroup
	for i := 0; i < numOfCrawlWorkers; i++ {
		wgCrawl.Add(1)
		go crawlConsumer(&wgCrawl, crawlCh, writeCh)
	}

	// Start write workers
	var wgWrite sync.WaitGroup
	for i := 0; i < numOfWriteWorkers; i++ {
		wgWrite.Add(1)
		go writeConsumer(&wgWrite, writeCh)
	}

	// Add URLs to crawl
	go func() {
		for _, url := range urls {
			crawlCh <- url
		}
		close(crawlCh)
	}()

	// Close writeCh when all crawl workers are done
	go func() {
		wgCrawl.Wait()
		close(writeCh)
	}()

	// Wait for all write workers to finish
	wgWrite.Wait()
}

func main() {
	urls := []string{
		"https://www.google.com",
	}

	CrawAndWrite(urls)

	// exit code 0
	os.Exit(0)
}

// normalizeFilename
func normalizeFilename(url string) string {
	re := regexp.MustCompile(`(?:https?://)?(?:www\.)?([^./]+)`)
	aa := re.FindStringSubmatch(url)
	var filename string
	if len(aa) > 1 {
		filename = aa[1]
	} else {
		filename = filename[:255]
	}
	return filename + ".html"
}
