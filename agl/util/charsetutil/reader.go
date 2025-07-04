package charsetutil

import (
	"bufio"
	"io"
	"regexp"

	"github.com/haorendashu/chardet"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/width"
)

var gbkDecoder = simplifiedchinese.GBK.NewDecoder()

// ChineseReader automatically read bytes encoded in different
// chinese format
type ChineseReader struct {
	r   *bufio.Reader
	dec io.Reader
}

// DetectCharset detects the charset.
func DetectCharset(b []byte) (string, error) {
	first := ""
	for _, v := range chardet.Possible(b) {
		if v == "utf-16be" || v == "utf-16le" {
			// Often confused with gbk.
			continue
		}
		if v == "utf-8" {
			return "utf-8", nil
		}
		if v == "hz-gb2312" {
			return "gbk", nil
		}
		if v == "gbk" {
			return "gbk", nil
		}
		if first == "" {
			first = v
		}
	}
	return first, nil
}

// NewChineseReader creates a new ChineseReader.
func NewChineseReader(r io.Reader) *ChineseReader {
	return &ChineseReader{r: bufio.NewReaderSize(r, 4096)}
}

// Read implements Reader, it ensure the output is utf8.
func (r *ChineseReader) Read(p []byte) (n int, err error) {
	if r.dec == nil {
		b, err := r.r.Peek(4096)
		if err != nil && err != io.EOF {
			return 0, err
		}
		detected, err := DetectCharset(b)
		if err != nil {
			return 0, err
		}
		switch detected {
		case "gbk", "hz-gb2312", "gb18030":
			// Just use gbk for all these cases, might not work,
			// if there is problem, provide a test case.
			r.dec = gbkDecoder.Reader(r.r)
		default:
			r.dec = r.r
		}
	}
	return r.dec.Read(p)
}

var hanPattern = regexp.MustCompile(`\p{Han}`)

// CleanTextWidth 全角转半角。
// 如果src中包含汉字则不做任何处理，否则返回转成半角后的字符。
func CleanTextWidth(src string) string {
	if hanPattern.MatchString(src) {
		return src
	}
	return width.Narrow.String(src)
}
