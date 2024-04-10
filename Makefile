tidy:
	@go mod tidy
	@go fmt

build:
	GOOS=linux GOARCH=amd64 go build -o birthdays

package: build
	zip birthdays.zip birthdays birthdays.csv

test:
	@go test -v ./...
