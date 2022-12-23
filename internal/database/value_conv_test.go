package database_test

import (
	"database/sql"
	"fmt"
	"github.com/kilianstallz/stage-sync/internal/database"
	"testing"
	"time"
)

func TestInt(t *testing.T) {
	t.Parallel()
	for scenario, fn := range map[string]func(t *testing.T){
		"returns int":        testValidInt,
		"returns float":      testValidFloat,
		"returns string":     testValidString,
		"returns bool true":  testValidBool,
		"returns bool false": testValidBoolFalse,
		"return time":        testValidTime,
		"return nil":         testInvalidInt,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

func testValidInt(t *testing.T) {
	nullInt := sql.NullInt64{Int64: 12, Valid: true}
	if fmt.Sprintf("%v", database.ConvertDbValue(nullInt)) != "12" {
		t.Errorf("nullInt.Int64 = %d, want 12", database.ConvertDbValue(nullInt))
	}
}

func testValidFloat(t *testing.T) {
	nullInt := sql.NullFloat64{Float64: 12.1, Valid: true}
	if fmt.Sprintf("%v", database.ConvertDbValue(nullInt)) != "12.1" {
		t.Errorf("nullInt.Float64 = %f, want 12.1", database.ConvertDbValue(nullInt))
	}
}

func testValidString(t *testing.T) {
	nullInt := sql.NullString{String: "a", Valid: true}
	if fmt.Sprintf("%v", database.ConvertDbValue(nullInt)) != "a" {
		t.Errorf("nullInt.String = %c, want 'a'", database.ConvertDbValue(nullInt))
	}
}

func testValidBool(t *testing.T) {
	nullInt := sql.NullBool{Bool: true, Valid: true}
	if fmt.Sprintf("%v", database.ConvertDbValue(nullInt)) != "true" {
		t.Errorf("nullInt.Bool = %c, want 'true'", database.ConvertDbValue(nullInt))
	}
}

func testValidBoolFalse(t *testing.T) {
	nullInt := sql.NullBool{Bool: false, Valid: true}
	if fmt.Sprintf("%v", database.ConvertDbValue(nullInt)) != "false" {
		t.Errorf("nullInt.Bool = %c, want 'false'", database.ConvertDbValue(nullInt))
	}
}

func testValidTime(t *testing.T) {
	ti := time.Now()
	nullInt := sql.NullTime{Time: ti, Valid: true}
	if fmt.Sprintf("%v", database.ConvertDbValue(nullInt)) != ti.String() {
		t.Errorf("nullInt.Time = %c, want 'a timestamp'", database.ConvertDbValue(nullInt))
	}
}

func testInvalidInt(t *testing.T) {
	nullInt := sql.NullInt64{Int64: 12, Valid: false}
	if database.ConvertDbValue(nullInt) != nil {
		t.Errorf("nullInt.Int64 = %d, want null", database.ConvertDbValue(nullInt))
	}
}
