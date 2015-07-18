#!/usr/bin/make -f

PROJECT=aws_deploy
SHELL=/bin/bash
MKDIR=mkdir
GIT=git
GO=go
RM=rm -rf

all: test

build: deps
	@echo "Building $(PROJECT)..."
	@$(GO) build -o bin/$(PROJECT)

deps:
	@echo "Downloading libraries..."
	@$(GO) get github.com/aws/aws-sdk-go/aws
	@$(GO) get github.com/aws/aws-sdk-go/service/ec2
	@$(GO) get github.com/aws/aws-sdk-go/service/elb
	@$(GO) get github.com/mitchellh/multistep
	@$(GO) get gopkg.in/alecthomas/kingpin.v1

dirs:
	@$(MKDIR) -p bin

test: build
	@echo "Run tests..."
	@$(GO) test ./...

clean:
	@echo "Cleanup binaries..."
	$(RM) bin
