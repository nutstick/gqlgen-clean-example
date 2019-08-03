# PID      = /tmp/awesome-golang-project.pid
# GO_FILES = $(wildcard *.go)
# APP      = ./app
# serve: restart
# 	@fswatch -o . | xargs -n1 -I{}  make restart || make kill

# kill:
# 	@kill `cat $(PID)` || true

# before:
# 	@echo "actually do nothing"
# $(APP): $(GO_FILES)
# 	@go build $? -o $@
# restart: kill before $(APP)
#         @app & echo $$! > $(PID)

# .PHONY: serve restart kill before # let's go to reserve rules names
SHELL := /bin/bash

# The name of the executable (default is current directory name)
TARGET := $(shell echo $${PWD\#\#*/})
.DEFAULT_GOAL: $(TARGET)

# These will be provided to the target
VERSION := 1.0.0
BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build clean install uninstall fmt simplify check run

all: check install

$(TARGET): $(SRC)
	@go build $(LDFLAGS) -o $(TARGET)

init:
	@glide install
	@go get -u github.com/jteeuwen/go-bindata/...
	@go get golang.org/x/tools/cmd/goimports

build: $(TARGET)
	@true

clean:
	@rm -f $(TARGET)

install:
	@go install $(LDFLAGS)

uninstall: clean
	@rm -f $$(which ${TARGET})

fmt:
	@gofmt -l -w $(SRC)

simplify:
	@gofmt -s -l -w $(SRC)

check:
	@test -z $(shell gofmt -l main.go | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done
	@go tool vet ${SRC}

run: install
	@$(TARGET)

dev:
	fresh

test:
	@go test -v ./...

gql:
	@go run scripts/gqlgen.go

deploy:
	@docker build -t gcr.io/nithi-project/nithi-backend-go .
	@gcloud docker -- push gcr.io/nithi-project/nithi-backend-go

deploy-debug:
	@docker build -t gcr.io/nithi-project/nithi-backend-go -f Dockerfile.debug .
	@gcloud docker -- push gcr.io/nithi-project/nithi-backend-go

mock:
	@sh script/mock.sh

gen: gql mock 