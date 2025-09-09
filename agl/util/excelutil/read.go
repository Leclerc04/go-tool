package excelutil

import (
	"fmt"

	"github.com/leclerc04/go-tool/agl/util/charsetutil"

	"bytes"

	"github.com/extrame/xls"
	"github.com/tealeg/xlsx"
)

const InvalidString = "A2_INVALID_STRING"

func ReadToString(b []byte, format string) ([][][]string, error) {
	var data [][][]string
	switch format {
	case "xlsx":
		f, err := xlsx.OpenBinary(b)
		if err != nil {
			return nil, err
		}
		for _, sheet := range f.Sheets {
			var sheetData [][]string
			for _, row := range sheet.Rows {
				var rowData []string
				for _, cell := range row.Cells {
					if v, err := cell.FormattedValue(); err == nil {
						rowData = append(rowData, charsetutil.CleanTextWidth(v))
					} else {
						rowData = append(rowData, InvalidString)
					}
				}
				sheetData = append(sheetData, rowData)
			}
			data = append(data, sheetData)
		}
	case "xls":
		r := bytes.NewReader(b)
		wb, err := xls.OpenReader(r, "utf-8")
		if err != nil {
			return data, err
		}
		// Read multiple sheets into the first element in data
		sheetData := wb.ReadAllCells(20000)
		for i, row := range sheetData {
			for ii, s := range row {
				sheetData[i][ii] = charsetutil.CleanTextWidth(s)
			}
		}
		data = append(data, sheetData)
	default:
		return nil, fmt.Errorf("%s is not a valid excel format", format)
	}
	return data, nil
}
