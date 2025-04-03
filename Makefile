# CLI Port Scanner Makefile
# This Makefile is used to build, test, and clean the CLI Port Scanner application.
# It includes targets for building the application, running tests, cleaning up build artifacts,
# formatting the source code, and running static analysis checks.

APP_NAME := cli-port-scanner
GO := go
PKG := ./...
PREFIX := [make]
BUILD_DIR := bin/

.DEFAULT_GOAL := build

fmt:
    @echo "$(PREFIX) Formatting source code..."
    @$(GO) fmt $(PKG)

vet: fmt
    @echo "$(PREFIX) Running vet to check code..."
    @$(GO) vet $(PKG)
	
build: vet
    @echo "$(PREFIX) Building $(APP_NAME) in $(BUILD_DIR)..."
    @mkdir -p $(BUILD_DIR)
    @$(GO) build -o $(BUILD_DIR)$(APP_NAME) main.go

test:
    @echo "$(PREFIX) Running tests..."
    @$(GO) test $(PKG) -v

clean:
    @echo "$(PREFIX) Cleaning build artifacts..."
    @rm -rfv $(BUILD_DIR)

check: fmt vet test
    @echo "$(PREFIX) Quality checks (format, vet, tests) complete!"

.PHONY: fmt vet build test clean check