OUT=out
BIN=${OUT}/bin
BINARY=${BIN}/mock-server
COVERAGE=${OUT}/coverage

.DEFAULT_GOAL := default
default: clean test-coverage build

tools:
	go get golang.org/x/tools/cmd/cover

test:
	go test -v ./...

test-coverage: tools
	rm -rf ${COVERAGE}
	mkdir -p ${COVERAGE}
	go test -coverprofile ${COVERAGE}/coverage.out -v ./...
	go tool cover -html ${COVERAGE}/coverage.out -o ${COVERAGE}/coverage.html

build: clean
	rm -rf ${BINARY}
	go build -o ${BINARY} .

clean:
	go clean
	
