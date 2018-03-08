package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeparateLinks(t *testing.T) {
	links := []string{"http://page.local/page1.html", "http://page.local/page2.local", "https://google.com"}
	expectedExternalLinkList := []string{"https://google.com"}
	expectedInternalLinkList := []string{"http://page.local/page1.html", "http://page.local/page2.local"}

	actualInternalLins, actualExternalLinks := separateLinks(links, "http://page.local")

	assert.Equal(t, expectedExternalLinkList, actualExternalLinks)
	assert.Equal(t, expectedInternalLinkList, actualInternalLins)
}

func TestIsArray(t *testing.T) {
	links := []string{"http://page.local/page1.html", "http://page.local/page2.html"}
	link := "http://page.local/page1.html"
	link2 := "http://page.local/page3.html"

	assert.True(t, isArray(links, link))
	assert.False(t, isArray(links, link2))
}
