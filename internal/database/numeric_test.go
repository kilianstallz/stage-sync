package database

import (
	numeric "github.com/jackc/pgtype/ext/shopspring-numeric"
	"testing"
)

func Test_Numeric(t *testing.T) {

	var nilNumeric numeric.Numeric
	err := nilNumeric.Set(nil)

	if err != nil {
		t.Errorf("Set(nil) error: %v", err)
	}
	v, err := nilNumeric.Value()
	if err != nil {
		t.Errorf("Value() error: %v", err)
	}

	if v != nil {
		t.Errorf("Value() = %v, want nil", v)
	}

	err = nilNumeric.Set(1.555555555555555)
	if err != nil {
		t.Errorf("Set(1.555555555555555) error: %v", err)
	}
	v, err = nilNumeric.Value()
	if err != nil {
		t.Errorf("Value() error: %v", err)
	}
	t.Log(v)
}
