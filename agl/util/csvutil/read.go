package csvutil

import (
	"encoding/csv"
	"io"

	"github.com/leclerc04/go-tool/agl/util/pipe"
)

// BOMutf8 is byte order mark for utf8 used for content detetion.
const BOMutf8 = "\xEF\xBB\xBF"

// WriteToReader writes csv data to io.Reader
func WriteToReader(thead []string, tbody [][]string) io.Reader {
	return pipe.WriterToReader(func(w io.Writer) error {
		cw := csv.NewWriter(w)
		// WriteAll will call flush.
		return cw.WriteAll(append([][]string{thead}, tbody...))
	})
}

func Write(w io.Writer, thead []string, tbody [][]string) error {
	_, err := w.Write([]byte(BOMutf8))
	if err != nil {
		return err
	}
	cw := csv.NewWriter(w)
	return cw.WriteAll(append([][]string{thead}, tbody...))
}
