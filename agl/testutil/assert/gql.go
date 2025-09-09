package assert

import (
	"testing"

	"github.com/leclerc04/go-tool/agl/util/simplejson"
	"github.com/stretchr/testify/assert"
)

// GQLError asserts the given graphql result contains some error message.
func GQLError(t *testing.T, r *simplejson.JSON, errMsg string) {
	if _, ok := r.CheckGet("errors"); !ok {
		assert.Fail(t, "no errors in result")
		return
	}
	assert.Contains(t, r.Get("errors", 0, "message").String(), errMsg)
}
