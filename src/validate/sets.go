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
