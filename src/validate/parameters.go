package validate

import (
	"errors"
)

/*
	Custom parameters implementation
*/

type MapParameters map[string]interface{}

type yamlExtra struct {
	value interface{}
}

func (p MapParameters) Get(name string) (interface{}, error) {

	value, found := p[name]
	// If no value is returned, try to grab from extra YAML
	if value == nil {
		value, found = p["data_"].(map[string]interface{})[name]
	}

	if !found {
		errorMessage := "No parameter '" + name + "' found."
		return nil, errors.New(errorMessage)
	}

	return value, nil
}
