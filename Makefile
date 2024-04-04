tidy:
	@go mod tidy
	@go fmt

build:
	@go build -o birthdays

package: build
	@zip birthdays.zip birthdays birthdays.csv

test:
	@go test -v ./...
