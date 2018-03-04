package main

// analyze given a url and a basurl, recoursively scans the page
import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// following all the links and fills the `visited` map
func analyze(url, baseurl string, visited *map[string]string, deep int) {
	fmt.Println("anaylysing......")
	page, err := parse(url)
	if err != nil {
		fmt.Printf("Error getting page %s %s\n", url, err)
		return
	}

	title := pageTitle(page)
	(*visited)[url] = title
	links := pageLinks(nil, page)
	assets := findAllStaticAssets(nil, page)
	deep++
	fmt.Printf("asset found:::::::::::::::::::::::::::::::::::: %+v\n", assets)

	for _, link := range links {
		//fmt.Println("checking the link if not visited", link)
		//fmt.Printf("link: %s and base url: %s", link, baseurl)
		//in this level base url
		link := baseurl + link
		fmt.Printf("go grab content from url:%s ", link)
		if (*visited)[link] == "" && strings.HasPrefix(link, baseurl) {
			fmt.Println("analyzing recursively ", link)
			if deep < 5 {
				analyze(link, baseurl, visited, deep)
			}
		}
	}
}

// pageTitle given a reference to a html.Node, scans it until it
// finds the title tag, and returns its value
func pageTitle(n *html.Node) string {
	var title string
	if n.Type == html.ElementNode && n.Data == "title" {
		return n.FirstChild.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title = pageTitle(c)
		if title != "" {
			break
		}
	}
	return title
}

// pageLinks will recursively scan a `html.Node` and will return
// a list of links found, with no duplicates
func pageLinks(links []string, n *html.Node) []string {
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
		links = pageLinks(links, c)
	}
	return links
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

// sliceContains returns true if `slice` contains `value`
func sliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// parse given a string pointing to a URL will fetch and parse it
// returning an html.Node pointer
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

// checkDuplicates scans the visited map for pages with duplicate titles
// and writes a report
func checkDuplicates(visited *map[string]string) {
	found := false
	uniques := map[string]string{}
	fmt.Printf("\nChecking duplicates..\n")
	for link, title := range *visited {
		if uniques[title] == "" {
			uniques[title] = link
		} else {
			found = true
			fmt.Printf("Duplicate title \"%s\" in %s but already found in %s\n", title, link, uniques[title])
		}
	}

	if !found {
		fmt.Println("No duplicates were found 😇")
	}
}
