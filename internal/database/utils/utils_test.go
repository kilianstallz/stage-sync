package utils_test

import (
	"stage-sync/internal/database/utils"
	"testing"
)

func TestArrayContains(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"zero length array valid":      testContainsZeroLength,
		"valid array contains":         testValidState,
		"valid array does not contain": testValidStateNotFound,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

func testContainsZeroLength(t *testing.T) {
	emptyArrInvalid := []string{}
	emptyArrValid := []string{}
	resFalse := utils.ArrayContains(emptyArrInvalid, "test")
	if resFalse {
		t.Errorf("Expected false got %v", resFalse)
	}

	resTrue := utils.ArrayContains(emptyArrValid, "")
	if resTrue {
		t.Errorf("Expected false got %v", resTrue)
	}
}

func testValidState(t *testing.T) {
	arr := []string{"id", "varlod", "af af"}
	resTrue := utils.ArrayContains(arr, "af af")
	if !resTrue {
		t.Errorf("Expected true got %v", resTrue)
	}
}

func testValidStateNotFound(t *testing.T) {
	arr := []string{"id", "varlod", "af af"}
	res := utils.ArrayContains(arr, "notContained")
	if res {
		t.Errorf("Expected false got %v", res)
	}
}
