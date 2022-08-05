tidy:
	go mod tidy

format:
	@GOFILES=$$(find . -name "*.go"); \
	for path in $${GOFILES}; do \
	 	gofmt -w -l $${path}; \
	  	goimports -w -l $${path}; \
	done;

check:
	golint -set_exit_status=1 ./...
	go vet ./...
	go test -test.timeout=20s -race -cpu=2 ./...

fix-golangci-lint:
	golangci-lint run -v --fix
