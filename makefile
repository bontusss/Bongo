## test: run all tests
test:
	@go test -v ./...


# cover: Opens coverage in browser
cover:
	@go test -coverprofile=cover.out ./... && go tool cover -html=cover.out


# coverage: displays test coverage
coverage:
	@go test -cover ./...

