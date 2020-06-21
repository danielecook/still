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
	// Add is_string()
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

func typeConvert(val string, NA_vals []string) interface{} {
	/*
		Automatically converts types
	*/
	for _, na := range NA_vals {
		if val == na {
			return nil
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

// TODO: Replace fname with iterator from excel, csv, etc.
func RunValidation(input string, schema schema.SchemaRules) bool {

	f, err := reader.NewReader(input, schema)
	utils.Check(err)

	colnames, err := f.ReadHeader()
	utils.Check(err)

	// Set all columns as valid to start
	ColumnStatus := make([]output.ValidCol, len(colnames))
	for k := range colnames {
		ColumnStatus[k].Name = colnames[k]
		ColumnStatus[k].IsValid = true
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

		// Set parameters
		parameters := make(map[string]interface{}, len(record))
		for idx := range record {
			parameters[colnames[idx]] = typeConvert(record[idx], schema.NA)
		}
		// Add additional parameters
		// Bools
		parameters["true"] = true
		parameters["false"] = false
		// ??? ADD parameters["NA"] = nil???

		for k, col := range schema.Columns {
			// Add in current column
			currentVar := typeConvert(record[indexOf(col.Name, colnames)], schema.NA)
			// TODO: Allow evaluation of NA values conditionally?
			if isNil(currentVar) {
				continue
			}
			parameters["current_var_"] = currentVar

			var funcSet = strings.Join(functionKeys(testFunctions), "|")

			var rule string
			// Allow for explcit references; They are removed (and added back later).
			explicitReplace, err := regexp.Compile(fmt.Sprintf("(%s)\\([ ]?%s[,]?", funcSet, col.Name))
			rule = explicitReplace.ReplaceAllString(col.Rule, "$1(")

			// Now add implicit argument back
			funcMatch, err := regexp.Compile(fmt.Sprintf("(%s)\\(", funcSet))
			utils.Check(err)

			rule = funcMatch.ReplaceAllString(rule, "$1(current_var_,")

			// If function takes single value, remove trailing comma
			rule = strings.Replace(rule, ",)", ")", -1)

			// The unique function needs a key; Add as implicit column and arguments
			uniqFunc, err := regexp.Compile("unique\\(([^)]+)")
			utils.Check(err)
			rule = uniqFunc.ReplaceAllString(rule, fmt.Sprintf("unique(\"%s:$1\",$1", col.Name))

			// TODO : Parse these just once!
			functions := combineFunctionSets(testFunctions, utilFunctions)
			expression, err := govaluate.NewEvaluableExpressionWithFunctions(rule, functions)
			if err != nil {
				log.Fatal(err)
			}
			result, err := expression.Evaluate(parameters)
			if err != nil {
				log.Fatal(err)
			}

			// Log results
			if result == false {
				ColumnStatus[k].IsValid = false
				ColumnStatus[k].NErrs++

				// Output log error
				fmt.Println(
					aurora.Sprintf("%s:%s[%d] %s -x-> '%s'",
						aurora.Red("Error"),
						aurora.Yellow(col.Name),
						f.Row(),
						aurora.Blue(col.Rule),
						currentVar))

			}
		}

	}

	output.PrintSummary(ColumnStatus)

	for _, i := range ColumnStatus {
		if i.IsValid == false {
			return false
		}
	}
	return true
}
