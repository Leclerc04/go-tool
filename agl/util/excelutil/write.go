package excelutil

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/tealeg/xlsx"
)

func WriteStringToXlsx(sheet [][]string) ([]byte, error) {
	return WriteSheets(nil, [][][]string{sheet})
}

func WriteSheets(sheetsName []string, sheets [][][]string) ([]byte, error) {
	file := xlsx.NewFile()
	for i, sheet := range sheets {
		var name string
		if len(sheetsName) <= i {
			name = fmt.Sprint("Sheet", i)
		} else {
			name = sheetsName[i]
		}
		s, err := file.AddSheet(name)
		if err != nil {
			return nil, err
		}

		for _, row := range sheet {
			r := s.AddRow()
			for _, cell := range row {
				c := r.AddCell()
				c.Value = cell
			}
		}
	}
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	if err := file.Write(w); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func WriteSheetsForCurriculumTable(sheetsName []string, sheets [][][]string) (*xlsx.File, error) {
	file := xlsx.NewFile()
	for i, sheet := range sheets {
		var name string
		if len(sheetsName) <= i {
			name = fmt.Sprint("Sheet", i)
		} else {
			name = sheetsName[i]
		}
		s, err := file.AddSheet(name)
		if err != nil {
			return nil, err
		}
		for _, row := range sheet {
			r := s.AddRow()
			for _, cell := range row {
				c := r.AddCell()
				style := c.GetStyle()
				style.ApplyAlignment = true
				style.Alignment.Vertical = "center"
				style.Alignment.Horizontal = "center"
				c.SetStyle(style)
				c.Value = cell
			}
		}
	}
	return file, nil
}
