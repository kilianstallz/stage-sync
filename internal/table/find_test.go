package table

import (
	"stage-sync/models"
	"testing"
)

func TestFindRow(t *testing.T) {
	rows := []models.Row{
		[]models.Column{
			{Name: "string", Type: "string", Value: "string"},
		},
	}

	rresult := FindRow(rows, []string{"string"}, rows[0])
	if rresult == nil {
		t.Errorf("should not be nil")
	}

	nilRes := FindRow([]models.Row{}, []string{}, models.Row{})
	if nilRes != nil {
		t.Errorf("Should be nil")
	}
}

func TestAllTrue(t *testing.T) {
	allT := []bool{true, true, true, true}
	allNT := []bool{true, true, true, false}
	allF := []bool{false, false, false, false}
	allFoT := []bool{false, false, false, true}
	if AllTrue(allT) == false {
		t.Errorf("Should be true")
	}
	if AllTrue(allNT) == true {
		t.Errorf("Should be false")
	}
	if AllTrue(allF) == true {
		t.Errorf("Should be false")
	}
	if AllTrue(allFoT) == true {
		t.Errorf("Should be false")
	}

}
