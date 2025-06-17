package main

import (
	"testing"
)

func TestBuildCells(t *testing.T) {
	var records [][]Cell
	records = append(records, []Cell{{Value: "foo", ValueType: "string"}})

	actual := jsonCellsToXmlCells(records)

	if actual[0][0].Text != "foo" {
		t.Fail()
	}
}
