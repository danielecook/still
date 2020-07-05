package reader

import (
	"fmt"

	csvReader "github.com/danielecook/still/src/reader/csv"
	excelReader "github.com/danielecook/still/src/reader/excel"
	"github.com/danielecook/still/src/schema"
	"github.com/danielecook/still/src/ui"
	"github.com/danielecook/still/src/utils"
	"github.com/gabriel-vasile/mimetype"
	"github.com/logrusorgru/aurora"
)

// https://github.com/lucmichalski/bigdata-stacks/blob/8fe9412b94dd5e0d20b14f29333037f6fc003757/refs/gleam/plugins/file/file_reader.go

type FileReader interface {
	Row() int
	Read() (row []string, err error)
	ReadHeader() (fieldNames []string, err error)
}

var plainCsv = map[string]rune{
	"text/tab-separated-values": '\t',
	"text/csv":                  ',',
	"text/plain; charset=utf-8": '\t',
}

var spreadSheets = []string{
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
}

func isPlain(s string) bool {
	if _, ok := plainCsv[s]; ok {
		return true
	}
	return false
}

func isSpreadsheet(s string) bool {
	for _, i := range spreadSheets {
		if i == s {
			return true
		}
	}
	return false
}

func inferDelim(s string) (rune, bool) {
	if isPlain(s) {
		return plainCsv[s], true
	}
	return ' ', false
}

func NewReader(fname string, schema schema.SchemaRules) (FileReader, error) {

	// Detect mimetype
	mime, err := mimetype.DetectFile(fname)
	utils.Check(err)

	// Set delimiter automatically for plaintext
	if schema.Separater == 0 {
		delim, ok := inferDelim(mime.String())
		if ok == true {
			var sepName string
			if delim == '\t' {
				sepName = "TAB"
			} else {
				sepName = ","
			}
			var sepMsg = "Separater inferred: [%v] ; Add %s to your schema to set manually"
			var sepDirective = fmt.Sprintf("@separater %s", sepName)
			ui.Warning(fmt.Sprintf(sepMsg, sepName, aurora.Green(sepDirective)))
			schema.Separater = delim
		}
	}
	switch {
	case isPlain(mime.String()):
		return csvReader.NewCSV(fname, schema), nil
	case isSpreadsheet(mime.String()):
		return excelReader.NewExcel(fname, schema), nil
	}
	return nil, fmt.Errorf("file type %s is not defined", fname)
}
