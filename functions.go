package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func getPageURLAndPageNumber() (string, int) {
	var numberOfPageToCrawl int
	var err error

	if len(os.Args) < 2 {
		fmt.Printf("Provide url to start crawling")
		os.Exit(1)
	}

	numberOfPageToCrawl = 10
	if len(os.Args) > 2 {
		if numberOfPageToCrawl, err = strconv.Atoi(os.Args[2]); err != nil {
			fmt.Println("Oops. wrong page number provided..")
			os.Exit(1)
		}
	}

	uri := os.Args[1]

	if uri == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	return uri, numberOfPageToCrawl
}
func crawl(url, baseURL string, visited *map[string]string, pageNumber int, wg *sync.WaitGroup, ch chan pageContent) {
	page, err := parse(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	title := getPageTitle(page)
	(*visited)[url] = title
	links := getPageLinks(nil, page)
	internalLinks, externalLinks := separateLinks(links, baseURL)
	assets := findAllStaticAssets(nil, page)

	item := pageContent{
		url,
		title,
		internalLinks,
		externalLinks,
		assets,
	}

	ch <- item

	for _, link := range internalLinks {
		if (*visited)[link] == "" && strings.HasPrefix(link, baseURL) {
			pageNumber = pageNumber - 1
			if pageNumber > 0 {
				fmt.Println("firing new crawl:", pageNumber)
				crawl(link, baseURL, visited, pageNumber, wg, ch)
			}
		}
	}
}

func getPageTitle(n *html.Node) string {
	var title string
	if n.Type == html.ElementNode && n.Data == "title" {
		return n.FirstChild.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title = getPageTitle(c)
		if title != "" {
			break
		}
	}

	return title
}

func getPageLinks(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				if !isArray(links, a.Val) {
					links = append(links, a.Val)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = getPageLinks(links, c)
	}

	return links
}

func isArray(items []string, value string) bool {
	for _, v := range items {
		if v == value {
			return true
		}
	}
	return false
}

func separateLinks(links []string, baseURL string) (internalLinks []string, externalLinks []string) {
	var currentURL string
	for _, link := range links {
		u, _ := url.Parse(link)
		if u.Host != "" {
			currentURL = u.Scheme + "://" + u.Host
			if strings.Contains(currentURL, baseURL) == false {
				externalLinks = append(externalLinks, link)
				continue
			}
		} else {
			link = baseURL + link
		}
		internalLinks = append(internalLinks, link)
	}

	return
}

func findAllStaticAssets(links []string, n *html.Node) []string {
	cssAsset := findCSSAssets(nil, n)
	links = append(links, cssAsset...)

	jsAssets := findJSAssets(nil, n)
	links = append(links, jsAssets...)

	return links
}

func findCSSAssets(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "link" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				if !isArray(links, a.Val) {
					links = append(links, a.Val)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = findCSSAssets(links, c)
	}

	return links
}

func findJSAssets(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "script" {
		for _, a := range n.Attr {
			if a.Key == "src" {
				if !isArray(links, a.Val) {
					links = append(links, a.Val)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = findJSAssets(links, c)
	}

	return links
}

func parse(url string) (*html.Node, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("can not crawl this page: %s", url)
	}

	respBody, err := html.Parse(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse page")
	}

	return respBody, err
}

func generatePageSiteMap(page pageContent) {
	fmt.Println("<page>")
	fmt.Printf(" <title>%s</title>\n", page.title)
	fmt.Printf(" <url>%s</url>\n", page.url)
	fmt.Println(" <internalLinks>")
	for _, link := range page.internalLinks {
		fmt.Printf("  <link>%s</link>\n", link)
	}
	fmt.Println(" </internalLinks>")

	fmt.Println(" <externalLinks>")
	for _, link := range page.externalLinks {
		fmt.Printf("  <link>%s</link>\n", link)
	}
	fmt.Println(" </externalLinks>")

	fmt.Println(" <assets>")
	for _, asset := range page.assets {
		fmt.Printf("  <asset>%s</asset>\n", asset)
	}
	fmt.Println(" </assets>")

	fmt.Println("</page>")
	fmt.Println()
}
