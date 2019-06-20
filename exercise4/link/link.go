package main

import (
	"bytes"
	"golang.org/x/net/html"
	"io"
	"strings"
)

// Link represents an HTML href and the text content.
type Link struct {
	Href string
	Text string
}

func isValid(link string) bool {
	if len(link) == 0 {
		return false
	}
	return true
}

func getNodeText(node *html.Node) string {
	var out bytes.Buffer
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.CommentNode || c.Type == html.DoctypeNode {
			continue
		} else if c.Type == html.TextNode {
			stripped := strings.Trim(c.Data, " \n\t")
			out.WriteString(stripped)
		} else {
			childText := getNodeText(c)
			if childText != "" {
				out.WriteString(" ")
				out.WriteString(childText)
			}
		}
	}
	return out.String()
}

func appendLinks(node *html.Node, links []Link) []Link {
	if node.Type == html.ElementNode && (node.Data == "a" || node.Data == "link") {
		for _, a := range node.Attr {
			if a.Key == "href" && isValid(a.Val) {
				text := getNodeText(node)
				links = append(links, Link{Href: a.Val, Text: text})
				break
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.CommentNode || c.Type == html.DoctypeNode {
			continue
		}
		links = appendLinks(c, links)
	}
	return links
}

// Scrap will return a list of links from an HTML document.
func Scrap(input io.Reader) ([]Link, error) {
	parsed, err := html.Parse(input)
	if err != nil {
		return nil, err
	}

	var links = []Link{}
	links = appendLinks(parsed, links)
	return links, nil
}
