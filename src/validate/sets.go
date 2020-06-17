package validate

func any(args ...interface{}) (interface{}, error) {
	// Checks for an element present in a set.
	for _, val := range args[1:] {
		if args[0] == val {
			return (bool)(true), nil
		}
	}
	return (bool)(false), nil
}

func rangeFunc(args ...interface{}) (interface{}, error) {
	var between bool
	switch val := args[0].(type) {
	case nil:
		return (bool)(true), nil
	case float64:
		between = args[1].(float64) <= val && val <= args[2].(float64)
	}
	val := args[0].(float64)
	between = args[1].(float64) <= val && val <= args[2].(float64)
	return (bool)(between), nil
}
