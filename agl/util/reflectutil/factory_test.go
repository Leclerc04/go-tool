package reflectutil_test

import (
	"fmt"
	"testing"

	"github.com/leclecr04/go-tool/agl/util/reflectutil"
)

type a struct{}

func TestFactory(t *testing.T) {
	{
		f := reflectutil.MakeFactory(&a{})
		sig := fmt.Sprintf("%T", f())
		if sig != "*reflectutil_test.a" {
			t.Error(sig)
		}
	}

	{
		f := reflectutil.MakeFactory(a{})
		sig := fmt.Sprintf("%T", f())
		if sig != "*reflectutil_test.a" {
			t.Error(sig)
		}
	}

}
