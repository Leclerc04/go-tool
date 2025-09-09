package html2text_test

import (
	"testing"

	. "github.com/leclerc04/go-tool/agl/util/html2text"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	assert.Equal(t, "foo bar", FromString("foo bar"))
	assert.Equal(t, "foo bar\nbar foo", FromString("foo bar<br />bar foo"))
	assert.Equal(t, "foo bar bar foo", FromString("foo bar <b>bar</b> foo"))
	assert.Equal(t, "foo bar\nbar\nfoo", FromString("foo bar <p>bar</p> foo"))
	assert.Equal(t, "foo bar [pic]foo", FromString("foo bar <img/>foo"))
	assert.Equal(t, "foo bar\nbaz foo", FromString("foo bar <ol><li>baz</li></ol> foo"))
}
