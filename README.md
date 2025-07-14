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

## License

This software is written by Florian Wilhelm and available under the MIT license (see `LICENSE` for details)
