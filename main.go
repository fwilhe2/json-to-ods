package main

import (
	"encoding/json"
	"flag"
	"os"
	"strings"

	rb "github.com/fwilhe2/rechenbrett"
)

type Cell struct {
	Value     string `json:"value"`
	ValueType string `json:"type"`
	Range     string `json:"range"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	flatPtr := flag.Bool("flat", false, "produce flat ods")
	inputFilePtr := flag.String("input", "spreadsheet.json", "input json file")
	outputFilePtr := flag.String("output", "spreadsheet.ods", "output (flat-)ods file")

	flag.Parse()

	dat, err := os.ReadFile(*inputFilePtr)
	check(err)

	var jsonCells [][]Cell
	var xmlCells [][]rb.Cell

	err = json.Unmarshal(dat, &jsonCells)
	check(err)

	for _, jsonRow := range jsonCells {
		var xmlRow []rb.Cell
		for _, jsonCell := range jsonRow {
			if len(jsonCell.Range) > 0 {
				xmlRow = append(xmlRow, rb.MakeRangeCell(jsonCell.Value, jsonCell.ValueType, jsonCell.Range))
			} else {
				xmlRow = append(xmlRow, rb.MakeCell(jsonCell.Value, jsonCell.ValueType))
			}
		}
		xmlCells = append(xmlCells, xmlRow)
	}

	spreadsheet := rb.MakeSpreadsheet(xmlCells)

	if *flatPtr {
		if strings.HasSuffix(*outputFilePtr, ".ods") {
			*outputFilePtr = strings.Replace(*outputFilePtr, ".ods", ".fods", -1)
		}
		os.WriteFile(*outputFilePtr, []byte(rb.MakeFlatOds(spreadsheet)), 0o644)
	} else {
		buff := rb.MakeOds(spreadsheet)

		archive, err := os.Create(*outputFilePtr)
		if err != nil {
			panic(err)
		}

		archive.Write(buff.Bytes())
		archive.Close()
	}
}
