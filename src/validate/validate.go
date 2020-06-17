package validate

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Knetic/govaluate"
	"github.com/danielecook/still/src/schema"
)

// TODO: Replace fname with iterator from excel, csv, etc.
func RunValidation(schema schema.SchemaRules, data string) bool {

	file, err := os.Open(data)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Test directives
	r := csv.NewReader(file)
	r.Comma = schema.Separater

	var i = 0
	var colnames = []string{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if i == 0 {
			colnames = record
			i = i + 1
			// TODO: Validate column order here
		}

		expression, err := govaluate.NewEvaluableExpression("10 > 0")
		result, err := expression.Evaluate(nil)
		fmt.Println(result)
		fmt.Printf("%#v\n", colnames)
		//fmt.Printf("%#v\n", record)
	}

	return true
}
