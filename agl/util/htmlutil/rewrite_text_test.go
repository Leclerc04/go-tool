package htmlutil_test

import (
	"context"
	"strings"
	"testing"

	"github.com/leclecr04/go-tool/agl/util/htmlutil"
	"github.com/stretchr/testify/assert"
)

func TestRewriteText(t *testing.T) {
	ctx := context.Background()
	ret, err := htmlutil.TransformHTMLText(
		ctx,
		`aaaa<a href="kkk">stanford</a><p>hello world stanford</p><a><p>hi there stanford</p></a><code>stanford</code>`,
		func(a string) (string, error) {
			a = strings.Replace(a, "stanford", `<a href="stanford_url">stanford</a>`, -1)
			return a, nil
		}, []string{"a", "code"})
	assert.NoError(t, err)
	assert.Equal(
		t, `aaaa<a href="kkk">stanford</a><p>hello world `+
			`<a href="stanford_url">stanford</a></p><a><p>hi there stanford</p></a><code>stanford</code>`,
		ret)
}
