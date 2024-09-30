.PHONY: run dev-run build build-dev test test-verbose submit stress-test

UNIT_TEST=go test -coverpkg=./... -coverprofile=coverage.out ./internal/.../...
INTEG_TEST=go test ./main_test.go
STRESS_TEST=go test ./stress_test.go
VERBOSE=-v

run:
	go run ./main.go

dev-run:
	env=dev go run ./main.go

build:
	go build -o bin/server ./main.go

build-dev:
	go build -o bin/server -ldflags "-X main.mode=dev" ./main.go

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

stress-test:
	$(STRESS_TEST) $(VERBOSE)

clean:
	- rm -rf images

submit:
	@echo "packing project for submission..."
	- zip -r receipt_uploader.zip . -x '.git/*' coverage.out *.zip bin/*
