# SPDX-FileCopyrightText: 2025 Florian Wilhelm
#
# SPDX-License-Identifier: MIT

all: format build test

format:
	gofumpt -w $$(find . -name '*.go')

build:
	go build -v ./...

test:
	go test -v ./...

install:
	sudo install json-to-ods /usr/local/bin

update:
	go get -u
	go mod tidy

demo:
	go run . -input sample.json -flat
	find samples -type f -name '*.json' ! -name 'hours-tracking.json' -exec go run . -input {} -flat -output {}.fods \;
	go run . -input samples/hours-tracking.json -flat -output samples/hours-tracking.json.fods \
		-table -header -structured-refs -banded -autofilter -totals "none,sum,sum,sum,sum,sum,sum"

clean:
	rm *ods
	rm *fods