package main

import (
	"strings"
	"testing"

	rb "github.com/fwilhe2/rechenbrett"
)

func TestBuildCells(t *testing.T) {
	var records [][]Cell
	records = append(records, []Cell{{Value: "foo", ValueType: "string"}})

	actual, err := jsonCellsToXmlCells(records)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if actual[0][0].Text != "foo" {
		t.Fail()
	}
}

func TestStyledCell(t *testing.T) {
	cells := [][]Cell{{{
		Value:     "foo",
		ValueType: "string",
		Style:     &Style{Bold: true, BackgroundColor: "#ffdc00"},
	}}}

	xml, err := jsonCellsToXmlCells(cells)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// MakeSpreadsheet is what materializes a styled cell into a named style.
	sheet, err := rb.MakeSpreadsheet(xml)
	if err != nil {
		t.Fatalf("MakeSpreadsheet: %v", err)
	}
	flat, err := rb.MakeFlatOds(sheet)
	if err != nil {
		t.Fatalf("MakeFlatOds: %v", err)
	}
	if !strings.Contains(flat, "#ffdc00") {
		t.Error("expected the styled cell's background color in the output")
	}
}

func TestRangeAndStyleAreMutuallyExclusive(t *testing.T) {
	cells := [][]Cell{{{
		Value:     "foo",
		ValueType: "string",
		Range:     "one",
		Style:     &Style{Bold: true},
	}}}

	if _, err := jsonCellsToXmlCells(cells); err == nil {
		t.Error("expected an error when a cell sets both range and style")
	}
}

func TestResolveColor(t *testing.T) {
	cases := map[string]string{
		"":        "",            // unset passes through
		"#00599d": "#00599d",     // literal hex passes through
		"navy":    rb.ColorNavy,  // palette name resolves
		"WHITE":   rb.ColorWhite, // case-insensitive
		"grey":    rb.ColorGray,  // british spelling aliases gray
	}
	for in, want := range cases {
		got, err := resolveColor(in)
		if err != nil {
			t.Errorf("resolveColor(%q) unexpected error: %v", in, err)
		}
		if got != want {
			t.Errorf("resolveColor(%q) = %q, want %q", in, got, want)
		}
	}

	if _, err := resolveColor("chartreuse"); err == nil {
		t.Error("expected an error for an unknown color name")
	}
}

func TestNamedColorInCell(t *testing.T) {
	cells := [][]Cell{{{
		Value:     "x",
		ValueType: "string",
		Style:     &Style{BackgroundColor: "navy"},
	}}}

	xml, err := jsonCellsToXmlCells(cells)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sheet, err := rb.MakeSpreadsheet(xml)
	if err != nil {
		t.Fatalf("MakeSpreadsheet: %v", err)
	}
	flat, err := rb.MakeFlatOds(sheet)
	if err != nil {
		t.Fatalf("MakeFlatOds: %v", err)
	}
	if !strings.Contains(flat, rb.ColorNavy) {
		t.Errorf("expected the resolved navy hex %q in the output", rb.ColorNavy)
	}
}

func TestParseTableStyle(t *testing.T) {
	cases := map[string]rb.TableStyle{
		"":      rb.TableStyleBlue,
		"blue":  rb.TableStyleBlue,
		"gray":  rb.TableStyleGray,
		"grey":  rb.TableStyleGray,
		"GREEN": rb.TableStyleGreen,
	}
	for in, want := range cases {
		got, err := parseTableStyle(in)
		if err != nil {
			t.Errorf("parseTableStyle(%q) unexpected error: %v", in, err)
		}
		if got != want {
			t.Errorf("parseTableStyle(%q) = %v, want %v", in, got, want)
		}
	}

	if _, err := parseTableStyle("chartreuse"); err == nil {
		t.Error("expected an error for an unknown table style")
	}
}

func TestParseTotals(t *testing.T) {
	totals, err := parseTotals("none, sum ,average,count,min,max")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []rb.TotalFunc{
		rb.TotalNone, rb.TotalSum, rb.TotalAverage, rb.TotalCount, rb.TotalMin, rb.TotalMax,
	}
	if len(totals) != len(want) {
		t.Fatalf("got %d totals, want %d", len(totals), len(want))
	}
	for i, w := range want {
		if totals[i].Func != w {
			t.Errorf("totals[%d].Func = %v, want %v", i, totals[i].Func, w)
		}
	}

	if got, _ := parseTotals("   "); got != nil {
		t.Error("expected nil totals for a blank spec")
	}
	if _, err := parseTotals("median"); err == nil {
		t.Error("expected an error for an unknown totals function")
	}
}
