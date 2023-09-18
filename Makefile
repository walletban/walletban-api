PROJECT_NAME = "walletban-api"
BASE=$(shell pwd)
BUILD_DIR=$(BASE)/bin
VERSION ?= "v1.0"
BUILD_DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT_SHA = $(shell git rev-parse --short HEAD)
LDFLAGS = -ldflags "-w -X main.version=${VERSION} -X main.buildDate=${BUILD_DATE} -X main.commit=${COMMIT_SHA}"
PACKAGE = $(shell go list -m)

.PHONY: clean
clean:
	@echo "> Cleaning Build targets"
	rm -rf bin

.PHONY: deps
deps:
	@echo "> Installing dependencies"
	@go mod tidy
	@go mod download

.PHONY: build
build: deps
	@echo "> Building obc-v0-platform backend Server Binary"
	go build ${LDFLAGS} -o ${BUILD_DIR}/${PROJECT_NAME}
	@echo "> Binary has been built successfully"


.PHONY: run
run: build
	@echo "> Running ${PROJECT_NAME}"
	${BUILD_DIR}/${PROJECT_NAME}