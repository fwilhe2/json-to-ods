package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	rb "github.com/fwilhe2/rechenbrett"
)

// Style is the optional per-cell appearance carried in the input JSON. Colors
// are hex strings such as "#ff0000"; Border is an ODF fo:border shorthand like
// "0.5pt solid #000000". It maps directly onto rechenbrett's CellStyle.
type Style struct {
	BackgroundColor string `json:"backgroundColor"`
	FontColor       string `json:"fontColor"`
	Bold            bool   `json:"bold"`
	Italic          bool   `json:"italic"`
	Border          string `json:"border"`
}

type Cell struct {
	Value     string `json:"value"`
	ValueType string `json:"type"`
	Range     string `json:"range"`
	Style     *Style `json:"style,omitempty"`
}

var version = "dev"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	flatPtr := flag.Bool("flat", false, "produce flat ods")
	inputFilePtr := flag.String("input", "spreadsheet.json", "input json file")
	outputFilePtr := flag.String("output", "spreadsheet.ods", "output (flat-)ods file")
	versionPtr := flag.Bool("version", false, "print version and exit")

	tablePtr := flag.Bool("table", false, "format the data as an Excel-style table (Format as Table)")
	headerPtr := flag.Bool("header", false, "treat the first row as a styled header (table mode)")
	bandedPtr := flag.Bool("banded", false, "shade alternating body rows (table mode)")
	autofilterPtr := flag.Bool("autofilter", false, "add AutoFilter dropdown buttons")
	structuredRefsPtr := flag.Bool("structured-refs", false, "create a named range per column, named after its header (table mode, requires -header)")
	tableStylePtr := flag.String("table-style", "blue", "table color theme: blue, gray or green (table mode)")
	tableNamePtr := flag.String("table-name", "", "name of the table / database range (table mode)")
	totalsPtr := flag.String("totals", "", "comma-separated totals row, one per column: none,sum,average,count,min,max (table mode)")

	flag.Parse()

	if *versionPtr {
		println("json-to-ods", version)
		println("https://github.com/fwilhe2/json-to-ods")
		println("Released under MIT License")
		println("Copyright (c) 2025 Florian Wilhelm")
		os.Exit(0)
	}

	dat, err := os.ReadFile(*inputFilePtr)
	check(err)

	var jsonCells [][]Cell

	err = json.Unmarshal(dat, &jsonCells)
	check(err)

	xmlCells, err := jsonCellsToXmlCells(jsonCells)
	check(err)

	var spreadsheet rb.Spreadsheet
	if *tablePtr {
		style, err := parseTableStyle(*tableStylePtr)
		check(err)
		totals, err := parseTotals(*totalsPtr)
		check(err)

		spreadsheet, err = rb.MakeTable(xmlCells, rb.TableOptions{
			Name:           *tableNamePtr,
			Header:         *headerPtr,
			AutoFilter:     *autofilterPtr,
			BandedRows:     *bandedPtr,
			Totals:         totals,
			StructuredRefs: *structuredRefsPtr,
			Style:          style,
		})
		check(err)
	} else {
		spreadsheet, err = rb.MakeSpreadsheet(xmlCells)
		check(err)
		if *autofilterPtr {
			spreadsheet = rb.EnableAutoFilter(spreadsheet)
		}
	}

	if *flatPtr {
		if strings.HasSuffix(*outputFilePtr, ".ods") {
			*outputFilePtr = strings.ReplaceAll(*outputFilePtr, ".ods", ".fods")
		}
		flat, err := rb.MakeFlatOds(spreadsheet)
		check(err)
		os.WriteFile(*outputFilePtr, []byte(flat), 0o644)
	} else {
		buff, err := rb.MakeOds(spreadsheet)
		check(err)

		archive, err := os.Create(*outputFilePtr)
		if err != nil {
			panic(err)
		}

		archive.Write(buff.Bytes())
		archive.Close()
	}
}

func jsonCellsToXmlCells(jsonCells [][]Cell) ([][]rb.Cell, error) {
	var xmlCells [][]rb.Cell

	for rowIdx, jsonRow := range jsonCells {
		var xmlRow []rb.Cell
		for colIdx, jsonCell := range jsonRow {
			cell, err := makeCell(jsonCell)
			if err != nil {
				return nil, fmt.Errorf("row %d, column %d: %w", rowIdx+1, colIdx+1, err)
			}
			xmlRow = append(xmlRow, cell)
		}
		xmlCells = append(xmlCells, xmlRow)
	}
	return xmlCells, nil
}

// makeCell builds a single rechenbrett cell from its JSON representation.
// rechenbrett offers no styled range cell, so range and style are mutually
// exclusive rather than silently dropping one of them.
func makeCell(c Cell) (rb.Cell, error) {
	hasRange := len(c.Range) > 0
	hasStyle := c.Style != nil

	switch {
	case hasRange && hasStyle:
		return rb.Cell{}, errors.New(`a cell may set either "range" or "style", not both`)
	case hasRange:
		return rb.MakeRangeCell(c.Value, c.ValueType, c.Range), nil
	case hasStyle:
		return rb.MakeStyledCell(c.Value, c.ValueType, rb.CellStyle{
			BackgroundColor: c.Style.BackgroundColor,
			FontColor:       c.Style.FontColor,
			Bold:            c.Style.Bold,
			Italic:          c.Style.Italic,
			Border:          c.Style.Border,
		}), nil
	default:
		return rb.MakeCell(c.Value, c.ValueType), nil
	}
}

func parseTableStyle(s string) (rb.TableStyle, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "blue":
		return rb.TableStyleBlue, nil
	case "gray", "grey":
		return rb.TableStyleGray, nil
	case "green":
		return rb.TableStyleGreen, nil
	default:
		return 0, fmt.Errorf("unknown table style %q (want blue, gray or green)", s)
	}
}

func parseTotals(s string) ([]rb.Total, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}

	var totals []rb.Total
	for part := range strings.SplitSeq(s, ",") {
		var f rb.TotalFunc
		switch strings.ToLower(strings.TrimSpace(part)) {
		case "", "none":
			f = rb.TotalNone
		case "sum":
			f = rb.TotalSum
		case "average", "avg":
			f = rb.TotalAverage
		case "count":
			f = rb.TotalCount
		case "min":
			f = rb.TotalMin
		case "max":
			f = rb.TotalMax
		default:
			return nil, fmt.Errorf("unknown totals function %q (want none, sum, average, count, min or max)", part)
		}
		totals = append(totals, rb.Total{Func: f})
	}
	return totals, nil
}
