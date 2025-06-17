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
	sudo install mkods /usr/local/bin

update:
	go get -u
	go mod tidy

demo:
	go run . -input sample.json -flat
	find samples -type f -name '*.json' -exec go run . -input {} -flat -output {}.fods \;

clean:
	rm *ods
	rm *fods