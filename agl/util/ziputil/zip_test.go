package ziputil_test

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/leclerc04/go-tool/agl/util/csvutil"
	"github.com/leclerc04/go-tool/agl/util/ziputil"
	"github.com/stretchr/testify/assert"
)

func TestZip(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	fm := os.FileMode(0640)
	assert.NoError(t, ziputil.Zip(&ziputil.File{
		Name:  "hello",
		IsDir: true,
		Files: []*ziputil.File{
			&ziputil.File{
				Name:     "love.txt",
				Body:     bytes.NewReader([]byte(`Cheers love!`)),
				FileMode: &fm,
			},
			&ziputil.File{
				Name:  "foo",
				IsDir: true,
				Files: []*ziputil.File{
					&ziputil.File{
						Name: "bar.csv",
						Body: csvutil.WriteToReader(
							[]string{"foo", "bar"},
							[][]string{[]string{"hahah", "hello"}}),
						FileMode: &fm,
					},
				},
			},
		},
	}, buf))
	bs := buf.Bytes()
	r := bytes.NewReader(bs)
	zr, err := zip.NewReader(r, r.Size())
	assert.NoError(t, err)

	readAll := func(t *testing.T, zipf *zip.File) string {
		rc, err := zipf.Open()
		assert.NoError(t, err)
		defer rc.Close()

		bs, err := ioutil.ReadAll(rc)
		assert.NoError(t, err)

		return string(bs)
	}

	assert.Equal(t, zr.File[0].Name, "hello/")
	assert.True(t, zr.File[0].Mode().IsDir())
	assert.Equal(t, os.FileMode(0666), zr.File[0].Mode().Perm())
	assert.Equal(t, os.FileMode(0666), zr.File[0].Mode().Perm())

	assert.Equal(t, zr.File[1].Name, "hello/love.txt")
	assert.True(t, zr.File[1].Mode().IsRegular())
	assert.Equal(t, os.FileMode(0640), zr.File[1].Mode().Perm())
	assert.Equal(t, `Cheers love!`, readAll(t, zr.File[1]))

	assert.Equal(t, zr.File[2].Name, "hello/foo/")
	assert.True(t, zr.File[2].Mode().IsDir())
	assert.Equal(t, os.FileMode(0666), zr.File[2].Mode().Perm())

	assert.Equal(t, zr.File[3].Name, "hello/foo/bar.csv")
	assert.True(t, zr.File[3].Mode().IsRegular())
	assert.Equal(t, os.FileMode(0640), zr.File[3].Mode().Perm())
	assert.Equal(t, "foo,bar\nhahah,hello\n", readAll(t, zr.File[3]))
}
