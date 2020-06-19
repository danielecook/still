package validate

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/danielecook/still/src/reader"
	"github.com/danielecook/still/src/schema"
	"github.com/danielecook/still/src/utils"
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
}

// Define functions
var testFunctions = map[string]govaluate.ExpressionFunction{
	// Test Functions
	"is":  is,
	"not": not,
	// Sets
	"any": any,
	// Strings
	"regex":     regex,
	"uppercase": uppercase,
	"lowercase": lowercase,
	"starts":    starts,
	"ends":      ends,
	"contains":  contains,
	"length":    length,
	// Numbers
	"range":       rangeFunc,
	"is_positive": isPositive,
	"is_negative": isNegative,
	// Types
	// TODO: ADD is_float()?
	// Add is_string()
	"is_numeric": isNumeric,
	"is_int":     isInt,
	"is_bool":    isBool,
	// Dates
	"is_date":         isDate,
	"is_date_relaxed": isDateRelaxed,
	"is_date_format":  isDateFormat,
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
func RunValidation(input string, schema schema.SchemaRules) bool {

	f, err := reader.NewReader(input, schema)
	utils.Check(err)

	colnames, err := f.ReadHeader()
	utils.Check(err)

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

		for _, col := range schema.Columns {
			// Add in current column
			currentVar := typeConvert(record[indexOf(col.Name, colnames)], schema.NA)
			fmt.Println(currentVar)
			// TODO: Allow evaluation of NA values conditionally?
			if isNil(currentVar) {
				continue
			}
			parameters["current_var_"] = currentVar

			var funcSet = strings.Join(functionKeys(testFunctions), "|")
			//var colSet = strings.Join(colnames, "|")
			var rule string
			// TODO: Need to test for explict column references and leave those intact;
			// if implicit, then replace; funcMatch.ReplaceAllStringFunc? Or
			// a more complex solution
			funcMatch, err := regexp.Compile(fmt.Sprintf("(%s)\\(", funcSet))
			if err != nil {
				log.Fatal(err)
			}
			rule = funcMatch.ReplaceAllString(col.Rule, "$1(current_var_,")
			// If function takes single value, remove trailing comma
			rule = strings.Replace(rule, ",)", ")", -1)
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
			fmt.Printf("%s: %s --> %v\n\n", col.Name, rule, result)
		}

	}

	return true
}
