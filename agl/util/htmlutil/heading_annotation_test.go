package htmlutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/leclerc04/go-tool/agl/util/htmlutil"
)

func TestHeadingAnnotation(t *testing.T) {
	data := `<div>hello</div><div><p>hey</p></div><div>_heading2: hello2 world</div><div>_heading3: hello3</div>`
	expected := "<div>hello</div><div><p>hey</p></div><h2>hello2 world</h2><h3>hello3</h3>"

	actual, err := AnnotationToHeading(data)
	assert.NoError(t, err)
	assert.EqualValues(t, expected, actual)

	actual, err = HeadingToAnnotation(actual)
	assert.NoError(t, err)
	assert.EqualValues(t, data, actual)
}
