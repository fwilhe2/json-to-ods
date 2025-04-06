package main

import (
	"encoding/json"
	rb "github.com/fwilhe2/rechenbrett"
	"os"
)

type Cell struct {
	Value     string `json:"value"`
	ValueType string `json:"type"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	dat, err := os.ReadFile("sample.json")
	check(err)

	var jsonCells [][]Cell
	var xmlCells [][]rb.Cell

	err = json.Unmarshal(dat, &jsonCells)
	check(err)

	for _, jsonRow := range jsonCells {
		var xmlRow []rb.Cell
		for _, jsonCell := range jsonRow {
			xmlRow = append(xmlRow, rb.MakeCell(jsonCell.Value, jsonCell.ValueType))
		}
		xmlCells = append(xmlCells, xmlRow)
	}

	println(jsonCells)
	spreadsheet := rb.MakeSpreadsheet(xmlCells)

	println(rb.MakeFlatOds(spreadsheet))
}
