package excelReader

import (
	"io"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/danielecook/still/src/schema"
	"github.com/danielecook/still/src/utils"
)

type excelReader struct {
	excelReader *excelize.File
	currentRow  int
	rows        [][]string
}

func NewExcel(fname string, sch schema.SchemaRules) *excelReader {
	file, err := excelize.OpenFile(fname)
	utils.Check(err)

	rows, err := file.GetRows(file.GetSheetList()[0])
	return &excelReader{
		excelReader: file,
		currentRow:  0,
		rows:        rows,
	}
}

func (r *excelReader) ReadHeader() (fieldNames []string, err error) {
	return r.rows[0], nil
}

func (r *excelReader) Read() (row []string, err error) {
	if (r.currentRow + 1) >= len(r.rows) {
		return []string{}, io.EOF
	}
	r.currentRow++
	return r.rows[r.currentRow], nil
}

func (r *excelReader) Row() int {
	return r.currentRow
}

func (r *excelReader) Reset() {
	r.currentRow = 0
}
