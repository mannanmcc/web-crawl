package main

import (
	"net/url"
	"sync"
)

func main() {
	wg := new(sync.WaitGroup)
	uri, numberOfPageToCrawl := getPageURLAndPageNumber()

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
