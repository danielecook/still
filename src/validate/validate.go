package validate

import (
	"fmt"
	"io"
	"log"
	"regexp"
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
	// Cumulative
	"last": last,
}

// Define functions
var testExpressions = map[string]govaluate.ExpressionFunction{
	// Test Functions
	"is":  is,
	"not": not,
	// Sets
	"any":            any,
	"identical":      identical,
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
	// Missing data
	"is_na":      isNA,
	"is_empty":   isEMPTY,
	"is_missing": isMissing,
}

var keyFunctions = []string{
	"unique",
	"identical",
	"last",
}

// RunValidation
func RunValidation(schema schema.SchemaRules, input string) bool {

	f, err := reader.NewReader(input, schema)
	utils.Check(err)

	colnames, err := f.ReadHeader()
	utils.Check(err)

	// Set all columns as valid or nil to start
	for idx := range schema.Columns {
		schema.Columns[idx].Status = 1 // 1=VALID
	}

	/*
		First compile expressions
	*/
	var rule string
	var expressions = make([]*govaluate.EvaluableExpression, len(schema.Columns))
	var funcSet = strings.Join(functionKeys(testExpressions), "|")
	var functions = combineFunctionSets(testExpressions, utilFunctions)

	//	Directive checks
	schema.IsOrdered(colnames)
	schema.IsFixed(colnames)

	for idx, col := range schema.Columns {

		// Allow for explcit references by removing them initialy
		explicitReplace, err := regexp.Compile(fmt.Sprintf("(%s)\\([ ]?%s[,]?", funcSet, col.Name))
		rule = explicitReplace.ReplaceAllString(col.Rule, "$1(")

		// Add implicit variables; Remove trailing commas
		funcMatch, err := regexp.Compile(fmt.Sprintf("(%s)\\(", funcSet))
		utils.Check(err)
		rule = funcMatch.ReplaceAllString(rule, "$1(current_var_,")
		rule = strings.Replace(rule, ",)", ")", -1)

		// Key functions require the variable name to create a key
		keyFunc, err := regexp.Compile(fmt.Sprintf("(%s)\\(([^)]+)", strings.Join(keyFunctions, "|")))
		utils.Check(err)
		rule = keyFunc.ReplaceAllString(rule, fmt.Sprintf("$1(\"%s:$2\",$2", col.Name))

		// If no expression is supplied set to true
		if rule == "" {
			rule = "true"
		}
		// Parse expressions
		expr, err := govaluate.NewEvaluableExpressionWithFunctions(rule, functions)
		if err != nil {
			log.Fatal(err)
		}
		expressions[idx] = expr
	}

	/*
		Then run the expressions on every row
	*/

	// Setup parameters with initial data
	parameters := make(MapParameters, len(colnames))
	parameters["true"] = true
	parameters["false"] = false
	parameters["data_"] = schema.YAMLData

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

		for idx := range record {
			parameters[colnames[idx]] = typeConvert(record[idx], schema.NA, schema.EMPTY)
		}

		for idx, col := range schema.Columns {
			colIndex := indexOf(col.Name, colnames)

			// Check to see if column exists
			if colIndex == -1 {
				schema.Columns[idx].Status = 3
				continue
			}

			// Add in current column
			currentVar := parameters[col.Name]
			parameters["current_var_"] = currentVar

			result, exprError := expressions[idx].Eval(parameters)

			// Column Summary
			if na, _ := isNA(currentVar); na.(bool) {
				schema.Columns[idx].NNA++
			} else if empty, _ := isEMPTY(currentVar); empty.(bool) {
				schema.Columns[idx].NEMPTY++
			} else if result == true {
				schema.Columns[idx].NVALID++
			}

			// Log results
			if result == false {
				schema.Columns[idx].Status = 2
				schema.Columns[idx].NErrs++

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
	output.PrintSummary(colnames, schema)

	// Fail if any single column fails or schema has failed
	for _, col := range schema.Columns {
		if col.Status == 2 {
			return false
		}
	}
	return true
}
