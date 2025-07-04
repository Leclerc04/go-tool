package htmlutil

import (
	"bytes"
	"context"
	"io"

	"github.com/leclecr04/go-tool/agl/base/trace"
	"github.com/leclecr04/go-tool/agl/util/strs"

	"golang.org/x/net/html"
)

// TransformHTMLText applies a transformation function to the text element of the HTML.
// It will not apply this function for text inside any of the skipTags.
func TransformHTMLText(ctx context.Context, input string, f func(string) (string, error), skipTags []string) (string, error) {
	z := html.NewTokenizer(bytes.NewBufferString(input))
	buf := bytes.Buffer{}

	skipDepth := 0
loop:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				break loop
			}
			// If the tokenization has problem, just ignore every thing
			// and return the raw input.
			trace.Printf(ctx, "tokenization has error: %v", z.Err())
			return input, nil
		case html.TextToken:
			if skipDepth > 0 {
				buf.Write(z.Raw())
				continue
			}
			newText, err := f(string(z.Text()))
			if err != nil {
				return "", err
			}
			buf.WriteString(newText)
		case html.StartTagToken:
			buf.Write(z.Raw())
			tagNameB, _ := z.TagName()
			tagName := string(tagNameB)
			if strs.InSlice(tagName, skipTags) {
				skipDepth++
			}
		case html.EndTagToken:
			buf.Write(z.Raw())
			tagNameB, _ := z.TagName()
			tagName := string(tagNameB)
			if strs.InSlice(tagName, skipTags) {
				skipDepth--
			}
		default:
			buf.Write(z.Raw())
		}
	}
	return buf.String(), nil
}
