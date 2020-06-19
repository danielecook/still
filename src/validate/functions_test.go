package validate

import "testing"

func TestRange(t *testing.T) {
	if val, err := rangeFunc(20.0, 10.0, 30.0); val != true {
		t.Errorf("range malfunction %s", err)
	}
	if val, err := rangeFunc(500.0, 10.0, 30.0); val != false {
		t.Errorf("range malfunction %s", err)
	}
}
