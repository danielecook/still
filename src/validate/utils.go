package validate

import (
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
)

func functionKeys(functions map[string]govaluate.ExpressionFunction) []string {
	evalFuncs := make([]string, len(functions))
	i := 0
	for k := range functions {
		evalFuncs[i] = k
		i++
	}
	return evalFuncs
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

// isNA is used to test
// for NA values internally
func isNA(args ...interface{}) bool {
	_, ok := args[0].(NA)
	if ok {
		return true
	}
	return false
}
