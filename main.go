package main

import (
	"net/url"
	"sync"
)

type counter struct {
	number int
}

func (c *counter) UpdateNumber() {
	c.number--
}

func (c *counter) getNumber() int {
	return c.number
}

func main() {
	wg := &sync.WaitGroup{}
	uri, numberOfPageToCrawl := getPageURLAndPageNumber()
	counter := &counter{number: numberOfPageToCrawl}

	u, _ := url.Parse(uri)
	baseURL := u.Scheme + "://" + u.Host

	visited := map[string]string{}
	pageContentChannel := make(chan pageContent)

	wg.Add(numberOfPageToCrawl + 1)
	go func() {
		defer wg.Done()
		counter.UpdateNumber()
		crawl(uri, baseURL, &visited, counter, wg, pageContentChannel)
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
