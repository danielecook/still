package validate

/*
	Missing Data

	NA - Known missing or "not applicable" data
	EMPTY - Unknown missing or currently unavailable data
*/
type NA string
type EMPTY string

func (c NA) String() string {
	return string(c)
}
func (c EMPTY) String() string {
	return string(c)
}

// isNA and isEMPTY are used to test for missing data internally
// for NA values internally
func isNA(args ...interface{}) (interface{}, error) {
	_, ok := args[0].(NA)
	if ok {
		return true, nil
	}
	return false, nil
}

func isEMPTY(args ...interface{}) (interface{}, error) {
	_, ok := args[0].(EMPTY)
	if ok {
		return true, nil
	}
	return false, nil
}

func isMissing(args ...interface{}) (interface{}, error) {
	na, _ := isNA(args[0])
	empty, _ := isEMPTY(args[0])
	return na.(bool) || empty.(bool), nil
}
