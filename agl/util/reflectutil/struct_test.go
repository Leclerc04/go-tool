package reflectutil_test

import (
	"testing"

	"github.com/leclerc04/go-tool/agl/util/reflectutil"
	"github.com/stretchr/testify/assert"
)

func TestFieldByName(t *testing.T) {
	st := struct {
		Boo  bool
		PBoo *bool
	}{}

	pBoo := reflectutil.FieldPtrByName(&st, "Boo").(*bool)
	*pBoo = true

	pPBoo := reflectutil.FieldPtrByName(&st, "PBoo").(**bool)
	*pPBoo = pBoo

	assert.True(t, st.Boo)
	assert.True(t, *st.PBoo)
}
