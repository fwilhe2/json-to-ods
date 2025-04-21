build:
	go build .

format:
	gofumpt -w main.go

update:
	go get -u
	go mod tidy

demo:
	go run . -input sample.json -flat

clean:
	rm *ods
	rm *fods