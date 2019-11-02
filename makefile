# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=hcunit
BINARY_DIR=build
BINARY_WIN=$(BINARY_NAME).exe
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_DARWIN=$(BINARY_NAME)_osx
BUILDTIME=$(shell date -u +%Y-%m-%d.%H:%M:%S)
BUMP_SEMVER_PATCH=$(shell git pull --tags >/dev/null && git tag -l | grep -v "-" | tail -1 | awk -F. '{print $$1"."$$2"."$$3+1}')
SHA_SHORT=$(shell git rev-parse --short HEAD)
SEMVER=$(BUMP_SEMVER_PATCH)-$(SHA_SHORT)
CLI_PATH=./cmd/hcunit

all: test build
build: build-darwin build-win build-linux 
test: gen unit e2e
unit: 
	$(GOTEST) ./pkg/... -v
e2e: 
	$(GOTEST) ./cmd/... -v
clean: 
	$(GOCLEAN)
	find . -name "*.test" | xargs rm 
	rm -fr $(BINARY_DIR)
gen:
	go generate ./...
build-darwin: 
	CGO_ENABLED=0 \
		GOOS=darwin \
		GOARCH=amd64 \
		$(GOBUILD) -ldflags "-X main.Buildtime=$(BUILDTIME) -X main.Version=$(SEMVER) -X main.Platform=OSX/amd64" -v -o $(BINARY_DIR)/$(BINARY_DARWIN) $(CLI_PATH) 
	chmod +x $(BINARY_DIR)/$(BINARY_DARWIN)
build-win:
	CGO_ENABLED=0 \
		GOOS=windows \
		GOARCH=amd64 \
		$(GOBUILD) -ldflags "-X main.Buildtime=$(BUILDTIME) -X main.Version=$(SEMVER) -X main.Platform=Windows/amd64"-v -o $(BINARY_DIR)/$(BINARY_WIN) $(CLI_PATH)
	chmod +x $(BINARY_DIR)/$(BINARY_WIN)
build-linux:
	CGO_ENABLED=0 \
		GOOS=linux \
		GOARCH=amd64 \
		$(GOBUILD) -ldflags "-X main.Buildtime=$(BUILDTIME) -X main.Version=$(SEMVER) -X main.Platform=Linux/amd64"-v -o $(BINARY_DIR)/$(BINARY_UNIX) $(CLI_PATH)
	chmod +x $(BINARY_DIR)/$(BINARY_UNIX)
release:
	./bin/create_new_release.sh

.PHONY: all test clean build
