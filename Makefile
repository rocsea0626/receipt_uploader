.PHONY: build build-dev test test-verbose submit

UNIT_TEST=go test -coverpkg=./... -coverprofile=coverage.out ./internal/utils/... ./internal/images/... ./internal/futils/... ./internal/handlers/...
INTEG_TEST=go test ./
VERBOSE=-v

build:
	go build -o bin/poker ./main.go

build-dev:
	go build -o bin/poker -ldflags "-X main.mode=dev" ./main.go

unit-test:
	$(UNIT_TEST)

integ-test:
	$(INTEG_TEST)

integ-test-verbose:
	$(INTEG_TEST) $(VERBOSE)

unit-test-verbose:
	$(UNIT_TEST) $(VERBOSE)

test:
	go clean -testcache
	@echo "running unit test..."
	make unit-test 
	@echo "running integration test..."
	make integ-test 

test-verbose:
	go clean -testcache
	@echo "running unit test in verbose mode..."
	make unit-test-verbose
	@echo "running integration test in verbose mode..."
	make integ-test-verbose

clean:
	- rm -rf bin

submit:
	@echo "packing project for submission..."
	- zip -r submit-larvis-poker.zip . -x '.git/*' coverage.out *.zip bin/* test_data_org.csv
