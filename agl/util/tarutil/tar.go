package tarutil

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/leclecr04/go-tool/agl/util/errs"
)

// UntarGZ untar the tar data to a target path.
func UntarGZ(reader io.Reader, target string) error {
	gzr, err := gzip.NewReader(reader)
	if err != nil {
		return errs.Wrap(err)
	}
	tarReader := tar.NewReader(gzr)

	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()

		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}

	return nil
}

// TarGZ arcives file from src file path.
func TarGZ(src string, writer io.Writer, skip func(path []string, content *[]byte) bool) error {
	// ensure the src actually exists before trying to tar it
	if _, err := os.Stat(src); err != nil {
		return errs.Wrap(err)
	}

	gzw := gzip.NewWriter(writer)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	wrapErr := func(err error, file string) error {
		return errs.Wrapf(err, "file: %s", file)
	}

	// walk path
	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		// return on any error
		if err != nil {
			return wrapErr(err, file)
		}

		ps := strings.Split(
			strings.TrimPrefix(
				strings.TrimPrefix(file, src),
				string(filepath.Separator),
			),
			string(filepath.Separator),
		)

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return wrapErr(err, file)
		}

		header.Name = strings.Join(ps, "/")

		if fi.Mode().IsDir() {
			err = tw.WriteHeader(header)
			return wrapErr(err, file)
		}

		bs, err := ioutil.ReadFile(file)
		if err != nil {
			return wrapErr(err, file)
		}
		if skip != nil && skip(ps, &bs) {
			return nil
		}

		if err := tw.WriteHeader(header); err != nil {
			return wrapErr(err, file)
		}
		_, err = tw.Write(bs)
		return err
	})
}
