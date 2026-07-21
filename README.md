# json-to-ods

CLI tool to convert json files into Open Document Spreadsheet files.

Powered by [rechenbrett](https://github.com/fwilhe2/rechenbrett).

Related: [csv-to-ods](https://github.com/fwilhe2/csv-to-ods)

## Motivation

This tool is helpful to transform existing data into a spreadsheet in cases where this might help with understanding or analyzing data.
The json format used is meant as an intermediate format that can be created by a script.

For example, if your data is in a complex xml or json structure, you might write a script in any language that is good for transforming data structures and leave the complexity of the Open Document Format to this too.

If you need more flexibility, feel free to build your own tool based on the [rechenbrett](https://github.com/fwilhe2/rechenbrett) library.

## Usage

Create a file called `input.json` with contents in a format like this:

```json
[
    [
        {
            "value": "foo",
            "type": "string"
        },
        {
            "value": "23.32",
            "type": "float",
            "range": "one"
        }
    ],
    [
        {
            "value": "2022-02-02",
            "type": "date"
        },
        {
            "value": "23.32",
            "type": "float",
            "range": "two"
        }
    ],
    [
        {
            "value": "SUM(B1:B2)",
            "type": "formula"
        },
        {
            "value": "AVERAGE(B1:B2)",
            "type": "formula"
        }
    ]
]
```

Use the cli like in this example:

```bash
json-to-ods -input input.json -flat -output output.fods
```

This will produce an Open Document Spreadsheet file.

If you omit the `flat` flag, a zipped Open Document Spreadsheet file will be created.
Use the `.ods` file extension in that case.

Check the [samples](./samples/) directory for more sample files.

## Cell types

The `type` of a cell may be one of:

| type | value format | notes |
| --- | --- | --- |
| `string` | any text | |
| `float` | `-?\d+(\.\d+)?` | |
| `percentage` | `-?\d+(\.\d+)?` | `0.15` renders as `15 %` |
| `currency` / `currency-eur` | `-?\d+(\.\d+)?` | euro |
| `currency-usd` | `-?\d+(\.\d+)?` | US dollar |
| `currency-gbp` | `-?\d+(\.\d+)?` | pound sterling |
| `date` | `YYYY-MM-DD` | |
| `time` | `HH:MM` or `HH:MM:SS` | |
| `formula` | e.g. `SUM(B1:B2)` | stored in OpenFormula notation |

## Cell styling

Any cell may carry an optional `style` object:

```json
{
    "value": "Total",
    "type": "string",
    "style": {
        "backgroundColor": "#ffdc00",
        "fontColor": "#111111",
        "bold": true,
        "italic": false,
        "border": "0.5pt solid #000000"
    }
}
```

A cell may use either `style` or `range`, but not both.

## Format as Table

Pass `-table` to mark the whole block as an Excel-style table ("Format as Table"),
with the following options:

| flag | effect |
| --- | --- |
| `-table` | enable table mode |
| `-header` | treat the first row as a styled header |
| `-banded` | shade alternating body rows |
| `-autofilter` | add AutoFilter dropdown buttons (also works without `-table`) |
| `-structured-refs` | create a named range per column, named after its header (requires `-header`) |
| `-table-style` | color theme: `blue` (default), `gray` or `green` |
| `-table-name` | name of the table / database range |
| `-totals` | comma-separated totals row, one entry per column: `none`, `sum`, `average`, `count`, `min`, `max` |

For example, to render a table with a header, banded rows, filter buttons,
column-named ranges and a summed totals row:

```bash
json-to-ods -input input.json -flat -output output.fods \
  -table -header -banded -autofilter -structured-refs \
  -table-style green -totals "none,sum,sum"
```

Totals use `SUBTOTAL`, so they respect the AutoFilter state (hidden rows are excluded).

## License

This software is written by Florian Wilhelm and available under the MIT license (see `LICENSE` for details)
