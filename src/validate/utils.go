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

func typeConvert(val string, NaVals []string, EmptyVals []string) interface{} {
	/*
		Automatically converts types
	*/
	for _, na := range NaVals {
		if val == na {
			return NA(na)
		}
	}

	for _, empty := range EmptyVals {
		if val == empty {
			return EMPTY(empty)
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
