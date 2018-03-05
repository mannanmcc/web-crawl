package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Provide url to start crawling")
		os.Exit(1)
	}

	uri := os.Args[1]

	if uri == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	//parse the url
	u, _ := url.Parse(uri)
	baseURL := u.Scheme + "://" + u.Host

	visited := map[string]string{}
	crawlingDepth := 3
	crawl(uri, baseURL, &visited, crawlingDepth)

	for link, title := range visited {
		fmt.Printf("%s -> %s\n", link, title)
	}
}
