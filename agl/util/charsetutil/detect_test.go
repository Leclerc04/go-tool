package charsetutil_test

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/leclecr04/go-tool/agl/testutil"
	"github.com/leclecr04/go-tool/agl/util/charsetutil"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		file    string
		charset string
	}{
		{
			file:    "test1.csv",
			charset: "gbk",
		},
		{
			file:    "test2.csv",
			charset: "utf-8",
		},
		{
			file:    "test3.csv",
			charset: "utf-8",
		},
	}

	for _, test := range tests {
		b, err := ioutil.ReadFile(path.Join(testutil.GetCurrentSourceFileDir(t), test.file))
		if err != nil {
			t.Error(err)
		}
		c, err := charsetutil.DetectCharset(b)
		if err != nil {
			t.Error(err)
		}
		if c != test.charset {
			t.Error(test, c)
		}
	}
}
