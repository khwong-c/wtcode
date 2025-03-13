.PHONY: all mod build

all: mod build

mod:
	go mod download

build:
	go build -o wtcode .

test:
	go test -v ./... -json | tee test-report.json