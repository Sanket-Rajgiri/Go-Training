build: 
	@go build -o binaries/grep main.go

test: 
	@go test 

coverage:
	@go test -coverprofile=coverage.out
	@go tool cover -func=coverage.out  