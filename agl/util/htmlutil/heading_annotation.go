package htmlutil

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/leclerc04/go-tool/agl/util/must"

	"golang.org/x/net/html"
)

var (
	annotationMatcher = regexp.MustCompile("_heading[2-6]:.*")
	headingMatcher    = regexp.MustCompile("h[2-6]$")
)

func AnnotationToHeading(s string) (string, error) {
	// returns the actrual text and heading size
	processAnnotation := func(text string) (string, int) {
		if !annotationMatcher.MatchString(text) {
			return text, 0
		}
		parts := strings.SplitN(text, ":", 2)
		i, err := strconv.Atoi(parts[0][len(parts[0])-1:])
		must.Must(err)
		return strings.TrimSpace(parts[1]), i
	}
	return process(s, func(n *html.Node) {
		if n.Type == html.ElementNode &&
			n.Data == "div" &&
			n.Attr == nil &&
			n.FirstChild.Type == html.TextNode {
			if t, h := processAnnotation(n.FirstChild.Data); h != 0 {
				n.FirstChild.Data = t
				n.Data = fmt.Sprintf("h%d", h)
			}
		}
	})
}

func HeadingToAnnotation(s string) (string, error) {
	// returns annotated text, and heading
	processHeading := func(text string, heading string) (string, string) {
		if !headingMatcher.MatchString(heading) {
			return text, heading
		}
		i, err := strconv.Atoi(heading[len(heading)-1:])
		must.Must(err)
		return fmt.Sprintf("_heading%d: %s", i, text), "div"
	}

	return process(s, func(n *html.Node) {
		if n.Type == html.ElementNode &&
			n.Attr == nil &&
			n.FirstChild != nil &&
			n.FirstChild.Type == html.TextNode {
			n.FirstChild.Data, n.Data = processHeading(n.FirstChild.Data, n.Data)
		}
	})
}

func process(s string, pf func(n *html.Node)) (string, error) {
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		return "", err
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		pf(n)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	buf := new(bytes.Buffer)
	err = html.Render(buf, doc)
	if err != nil {
		return "", err
	}
	s = buf.String()
	return strings.TrimSuffix(strings.TrimPrefix(s, "<html><head></head><body>"), "</body></html>"), nil
}
