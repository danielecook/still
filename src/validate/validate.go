package validate

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/danielecook/still/src/schema"
)

// Define functions
var functions = map[string]govaluate.ExpressionFunction{
	"strlen": func(args ...interface{}) (interface{}, error) {
		length := len(args[0].(string))
		return (float64)(length), nil
	},
	"in_range": func(args ...interface{}) (interface{}, error) {
		fmt.Printf("%#v", args)
		return (bool)(true), nil
	},
}

func indexOf(word string, data []string) int {
	for k, v := range data {
		if word == v {
			return k
		}
	}
	return -1
}

func typeConvert(val string) interface{} {
	/*
		Automatically converts types
	*/

	// Is it a Bool?
	if strings.ToUpper(val) == "TRUE" {
		return true
	}
	if strings.ToUpper(val) == "FALSE" {
		return false
	}

	// Is it an Integer?
	valInt, err := strconv.Atoi(val)
	if err == nil {
		return valInt
	}

	// Is it a float?
	valFloat, err := strconv.ParseFloat(val, 64)
	if err == nil {
		return valFloat
	}

	// string
	return val
}

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
	stopRead := false
	for ok := true; ok; ok = (stopRead == false) {
		record, readErr := r.Read()
		if readErr == io.EOF {
			stopRead = true
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		i = i + 1
		if i == 1 {
			colnames = record
			continue
			// TODO: Validate column order here
		}

		// Set parameters
		parameters := make(map[string]interface{}, len(record))
		for idx := range record {
			parameters[colnames[idx]] = typeConvert(record[idx])
		}

		for _, col := range schema.Columns {
			// Add in current column
			parameters["current_var_"] = typeConvert(record[indexOf(col.Name, colnames)])

			evalFuncs := make([]string, len(functions))
			for k, _ := range functions {
				evalFuncs = append(evalFuncs, k)
			}
			var rule string
			for _, function := range evalFuncs {
				rule = strings.Replace(col.Rule,
					fmt.Sprintf("%s(", function),
					fmt.Sprintf("%s(current_var_,", function), 1000)
			}
			fmt.Println(rule)

			// TODO : Parse these just once!
			expression, err := govaluate.NewEvaluableExpressionWithFunctions(rule, functions)
			if err != nil {
				log.Fatal(err)
			}
			result, err := expression.Evaluate(parameters)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%v\n", result)
		}

	}

	return true
}
