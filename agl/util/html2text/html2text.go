package html2text

import (
	"bytes"
	"io"
	"log"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var (
	cleanNewLine = regexp.MustCompile("\\s*\n\\s*")
	cleanSpace   = regexp.MustCompile(" ")
)

func extract(node *html.Node, buff *bytes.Buffer) {
	inNewLine := false
	if node.Type == html.TextNode {
		buff.WriteString(node.Data)
	} else if node.Type == html.ElementNode && node.Data == "br" {
		buff.WriteString("\n")
	} else if node.Type == html.ElementNode {
		if node.Data == "p" {
			inNewLine = true
			buff.WriteString("\n")
		}
		if node.Data == "script" {
			return
		}
		if node.Data == "img" {
			buff.WriteString("[pic]")
			return
		}
		if node.Data == "li" {
			buff.WriteString("\n")
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		extract(c, buff)
	}
	if inNewLine {
		buff.WriteString("\n")
	}
}

func FromString(s string) string {
	return FromReader(strings.NewReader(s))
}

func FromReader(reader io.Reader) string {
	var buffer bytes.Buffer
	doc, err := html.Parse(reader)
	if err != nil {
		log.Fatal(err)
	}
	extract(doc, &buffer)
	s := strings.Trim(buffer.String(), "\n ")
	return cleanSpace.ReplaceAllString(cleanNewLine.ReplaceAllString(s, "\n"), " ")
}
