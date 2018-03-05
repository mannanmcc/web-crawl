package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func crawl(url, baseURL string, visited *map[string]string, pageNumber int, wg *sync.WaitGroup, ch chan pageContent) {
	page, err := parse(url)

	if err != nil {
		fmt.Printf("Error getting page %s %s\n", url, err)
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
				if !sliceContains(links, a.Val) {
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
				if !sliceContains(links, a.Val) {
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
				if !sliceContains(links, a.Val) {
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

func sliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func parse(url string) (*html.Node, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Cannot get page")
	}
	b, err := html.Parse(r.Body)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse page")
	}

	return b, err
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
