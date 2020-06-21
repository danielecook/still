package csvReader

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/danielecook/still/src/schema"
)

type CsvFileReader struct {
	csvReader *csv.Reader
	row       int
}

func NewCSV(fname string, sch schema.SchemaRules) *CsvFileReader {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(file)
	r.Comma = sch.Separater
	return &CsvFileReader{
		csvReader: r,
		row:       0,
	}
}

func (r *CsvFileReader) ReadHeader() (fieldNames []string, err error) {
	r.row++
	return r.csvReader.Read()
}

func (r *CsvFileReader) Read() (row []string, err error) {
	r.row++
	return r.csvReader.Read()
}

func (r *CsvFileReader) Row() int {
	return r.row
}
