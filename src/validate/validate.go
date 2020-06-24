package validate

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/danielecook/still/src/output"
	"github.com/danielecook/still/src/reader"
	"github.com/danielecook/still/src/schema"
	"github.com/danielecook/still/src/utils"
	"github.com/logrusorgru/aurora"
)

func isNil(val interface{}) bool {
	return nil == val
}

var utilFunctions = map[string]govaluate.ExpressionFunction{
	// Utility Functions
	"str_len": strLen,
	"max":     maxFunc,
	"min":     minFunc,
	"if_else": ifElse,
	"count":   countFunc,
	"print":   print,
	// strings
	"to_upper": toUpper,
	"to_lower": toLower,
	"replace":  replace,
}

// Define functions
var testFunctions = map[string]govaluate.ExpressionFunction{
	// Test Functions
	"is":  is,
	"not": not,
	// Sets
	"any":            any,
	"last":           last,
	"unique":         uniqueFunc,
	"is_subset_list": isSubsetList,
	// Strings
	"regex":     regex,
	"uppercase": uppercase,
	"lowercase": lowercase,
	"starts":    starts,
	"ends":      ends,
	"contains":  contains,
	"length":    length,
	"is_url":    isURL,
	// Numbers
	"range":       rangeFunc,
	"is_positive": isPositive,
	"is_negative": isNegative,
	// Types
	"is_string":  isString,
	"is_numeric": isNumeric,
	"is_int":     isInt,
	"is_bool":    isBool,
	// Dates
	"is_date":         isDate,
	"is_date_relaxed": isDateRelaxed,
	"is_date_format":  isDateFormat,
	// Files
	"file_exists":   fileExists,
	"file_min_size": fileMinSize,
	"file_max_size": fileMaxSize,
	"mimetype":      mimeTypeIs,
}

var keyFunctions = []string{
	"unique",
	"last",
}

func combineFunctionSets(ms ...map[string]govaluate.ExpressionFunction) map[string]govaluate.ExpressionFunction {
	res := map[string]govaluate.ExpressionFunction{}
	for _, m := range ms {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
}

func functionKeys(functions map[string]govaluate.ExpressionFunction) []string {
	evalFuncs := make([]string, len(functions))
	i := 0
	for k := range functions {
		evalFuncs[i] = k
		i++
	}
	return evalFuncs
}

func indexOf(word string, data []string) int {
	for k, v := range data {
		if word == v {
			return k
		}
	}
	return -1
}

type NA string

func (c NA) String() string {
	fmt.Println("Executing String() for NA!")
	return string(c)
}

func typeConvert(val string, NA_vals []string) interface{} {
	/*
		Automatically converts types
	*/
	for _, na := range NA_vals {
		if val == na {
			return NA(na)
		}
	}

	// bool
	if strings.ToUpper(val) == "TRUE" {
		return true
	}
	if strings.ToUpper(val) == "FALSE" {
		return false
	}

	// integer
	valInt, err := strconv.Atoi(val)
	if err == nil {
		return valInt
	}

	// float
	valFloat, err := strconv.ParseFloat(val, 64)
	if err == nil {
		return valFloat
	}

	// string
	return val
}

// RunValidation
func RunValidation(input string, schema schema.SchemaRules) bool {

	f, err := reader.NewReader(input, schema)
	utils.Check(err)

	colnames, err := f.ReadHeader()
	utils.Check(err)

	// Set all columns as valid or nil to start
	ColumnStatus := make([]output.ValidCol, len(colnames))
	for k := range colnames {
		ColumnStatus[k].Name = colnames[k]
		ColumnStatus[k].IsValid = 0
		for _, schemaCol := range schema.Columns {
			if schemaCol.Name == ColumnStatus[k].Name {
				ColumnStatus[k].IsValid = 1
			}
		}
	}

	stopRead := false
	for ok := true; ok; ok = (stopRead == false) {
		record, readErr := f.Read()
		if readErr == io.EOF {
			stopRead = true
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		/*
			Once expressions are compiled run on every row.
		*/
		for idx, col := range schema.Columns {

			// Add in current column
			currentVar := typeConvert(record[indexOf(col.Name, colnames)], schema.NA)

			// Set parameters
			parameters := make(MapParameters, len(record))
			for idx := range record {
				parameters[colnames[idx]] = typeConvert(record[idx], schema.NA)
			}

			// Add additional parameters
			// Bools
			parameters["true"] = true
			parameters["false"] = false

			// set parameters
			parameters["current_var_"] = currentVar
			parameters["data_"] = schema.YAMLData

			// Key functions require the variable name to create a key
			keyFunc, err := regexp.Compile(fmt.Sprintf("(%s)\\(([^)]+)", strings.Join(keyFunctions, "|")))
			utils.Check(err)

			result, exprError := expression.Eval(parameters)

			// Log results
			if result == false {
				colIndex := indexOf(col.Name, colnames)
				ColumnStatus[colIndex].IsValid = 2
				ColumnStatus[colIndex].NErrs++

				var infoLine string
				if exprError != nil {
					infoLine = aurora.Sprintf("%s (%s)", col.Rule, exprError)
				} else {
					infoLine = aurora.Sprintf("%s", col.Rule)
				}

				// Output log error
				fmt.Println(
					aurora.Sprintf("%5s:%-20s\t%5s\t%-15v\t→\t%s",
						aurora.Red("Error"),
						aurora.Yellow(col.Name),
						fmt.Sprintf("[%d]", f.Row()),
						aurora.Yellow(currentVar),
						aurora.Blue(infoLine)))
			}
		}

	}

	output.PrintSummary(ColumnStatus)

	// Fail if any single column fails
	for _, i := range ColumnStatus {
		if i.IsValid == 2 {
			return false
		}
	}
	return true
}
