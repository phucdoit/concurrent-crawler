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

func crawl(wg *sync.WaitGroup, crawlCh <-chan string, writeCh chan<- CrawlData) {
	defer wg.Done()
	for v := range crawlCh {
		resp, err := http.Get(v)
		if err != nil {
			writeCh <- CrawlData{Url: v, Err: err}
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			writeCh <- CrawlData{Url: v, Err: err}
			return
		}
		writeCh <- CrawlData{Url: v, Data: string(body)}
	}
}

func write(wg *sync.WaitGroup, writeCh <-chan CrawlData) {
	defer wg.Done()
	for v := range writeCh {
		if v.Err != nil {
			log.Println("URL:", v.Url, "Error:", v.Err)
			continue
		}

		if err := os.MkdirAll("output", os.ModePerm); err != nil {
			log.Fatal(err)
		}

		fileName := "output/" + normalizeFilename(v.Url)
		f, err := os.Create(fileName)
		if err != nil {
			log.Println("Error creating file", fileName, ":", err)
			continue
		}

		_, err = io.WriteString(f, v.Data)
		if err != nil {
			log.Println("Error writing to file", fileName, ":", err)
		}
		f.Close()
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
		go crawl(&wgCrawl, crawlCh, writeCh)
	}

	// Start write workers
	var wgWrite sync.WaitGroup
	for i := 0; i < numOfWriteWorkers; i++ {
		wgWrite.Add(1)
		go write(&wgWrite, writeCh)
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
