package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	wg := new(sync.WaitGroup)
	var numberOfPageToCrawl int
	var err error

	if len(os.Args) < 2 {
		fmt.Printf("Provide url to start crawling")
		os.Exit(1)
	}

	numberOfPageToCrawl = 2
	if len(os.Args) > 2 {
		if numberOfPageToCrawl, err = strconv.Atoi(os.Args[2]); err != nil {
			fmt.Println("Oops. wrong data provided..")
			os.Exit(1)
		}
	}

	uri := os.Args[1]

	if uri == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	u, _ := url.Parse(uri)
	baseURL := u.Scheme + "://" + u.Host

	visited := map[string]string{}

	pageContentChannel := make(chan pageContent)
	wg.Add(1)
	go func() {
		defer wg.Done()
		crawl(uri, baseURL, &visited, numberOfPageToCrawl, wg, pageContentChannel)
	}()

	go func() {
		wg.Wait()
		close(pageContentChannel)
	}()

	for page := range pageContentChannel {
		generatePageSiteMap(page)
	}
	wg.Wait()
}
