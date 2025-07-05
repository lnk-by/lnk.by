# Root Makefile for building, testing, and deploying Go Lambda functions

BINS_DIR=bin
GOOS=linux
GOARCH=amd64

LAMBDA_TARGETS = \
	aws/campaign/create \
	aws/campaign/delete \
	aws/campaign/list \
	aws/campaign/retrieve \
	aws/campaign/update \
	aws/customer/create \
	aws/customer/delete \
	aws/customer/list \
	aws/customer/retrieve \
	aws/customer/update \
	aws/organization/create \
	aws/organization/delete \
	aws/organization/list \
	aws/organization/retrieve \
	aws/organization/update \
	aws/redirect \
	aws/shorturl/create \
	aws/shorturl/delete \
	aws/shorturl/list \
	aws/shorturl/retrieve \
	aws/shorturl/update

SUBMODULES = \
	shared \
	aws \
	server

all: test build

test:
	@echo "mode: set" > coverage.out
	@for module in $(SUBMODULES); do \
		echo "Running tests in $$module..."; \
		go test -coverprofile=coverage_tmp.out ./$$module/... || true; \
		if [ -f coverage_tmp.out ]; then \
			tail -n +2 coverage_tmp.out >> coverage.out; \
			rm coverage_tmp.out; \
		fi \
	done
	@echo "Generating coverage report..."
	@go tool cover -func=coverage.out || true

build:
	mkdir -p $(BINS_DIR)
	@for target in $(LAMBDA_TARGETS); do \
		out_name=$$(echo $$target | sed 's|/|_|g'); \
		echo "Building $$target -> $(BINS_DIR)/$$out_name"; \
		GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINS_DIR)/$$out_name ./$$target; \
	done
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINS_DIR)/server server/main.go

clean:
	rm -rf $(BINS_DIR) coverage.out

.PHONY: all test build clean

