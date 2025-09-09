package ziputil

import (
	"archive/zip"
	"io"
	"os"
	"strings"
	"time"

	"github.com/leclerc04/go-tool/agl/util/errs"

	"github.com/leclerc04/go-tool/agl/util/timeutil"
)

// File represents file/directory to be zipped.
type File struct {
	Name     string
	IsDir    bool
	Files    []*File
	Body     io.Reader
	FileMode *os.FileMode // Optional.
}

// Zip compresses f and write to w.
func Zip(f *File, w io.Writer) (err error) {
	zw := zip.NewWriter(w)
	defer func() {
		err2 := zw.Flush()
		if err == nil {
			err = err2
		}
		err3 := zw.Close()
		if err == nil {
			err = err3
		}
	}()
	now := timeutil.Now()

	return zipInternal([]string{f.Name}, f, now, zw)
}

func zipInternal(path []string, f *File, now time.Time, w *zip.Writer) error {
	if f.IsDir {
		if f.Body != nil {
			return errs.Newf("Directory cannot contain body. path: %s", path)
		}
		h := &zip.FileHeader{
			Name:   strings.Join(path, "/") + "/",
			Method: zip.Store,
		}
		h.Modified = now
		if f.FileMode != nil {
			if !f.FileMode.IsDir() {
				return errs.Newf("FileMode %s not a directory. Path %s", f.FileMode.String(), path)
			}
			h.SetMode(*f.FileMode)
		}

		_, err := w.CreateHeader(h)
		if err != nil {
			return err
		}
		for _, c := range f.Files {
			err := zipInternal(append(path, c.Name), c, now, w)
			if err != nil {
				return err
			}
		}
		return nil
	}
	if len(f.Files) != 0 {
		return errs.Newf("Not directory found %v files, path: %s", len(f.Files), path)
	}
	h := &zip.FileHeader{
		Name:   strings.Join(path, "/"),
		Method: zip.Deflate,
	}
	h.Modified = now
	if f.FileMode != nil {
		if !f.FileMode.IsRegular() {
			return errs.Newf("FileMode %s is not regular. Path: %s", f.FileMode.String(), path)
		}
		h.SetMode(*f.FileMode)
	}
	ff, err := w.CreateHeader(h)
	if err != nil {
		return err
	}
	_, err = io.Copy(ff, f.Body)
	return err
}
