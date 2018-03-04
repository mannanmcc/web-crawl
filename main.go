package main

import (
	"flag"
	"fmt"
	"os"
)

type node struct {
	Type                    NodeType
	Data                    string
	Attr                    []Attribute
	FirstChild, NextSibling *node
}

type pageContent struct {
	title string
	links []string
}

type NodeType int32

const (
	ErrorNode NodeType = iota
	TextNode
	DocumentNode
	ElementNode
	CommentNode
	DoctypeNode
)

type Attribute struct {
	Key, Val string
}

func main() {
	var url string
	var dup bool
	flag.StringVar(&url, "url", "http://page.local", "the url to parse")
	flag.BoolVar(&dup, "dup", false, "if set, check for duplicates")
	flag.Parse()

	if url == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	visited := map[string]string{}

	deep := 1
	analyze(url, url, &visited, deep)
	for link, title := range visited {
		fmt.Printf("%s -> %s\n", link, title)
	}

	if dup {
		checkDuplicates(&visited)
	}
}
