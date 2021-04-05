SRC=$(shell find . -name "*.go")

# project information
PROJECTNAME=$(shell basename "$(PWD)")
GOPKGNAME = github.com/hamed-yousefi/athenz-agent
PKG_DATE=$(shell date '+%Y-%m-%dT%H:%M:%S')
SRC=cmd/tools/athenz-agent.go
GOBASE=$(PWD)
GOBIN=$(GOBASE)/bin
GOPATH=$(GOBASE)/vendor

# build information
BUILDPATH=/home/athenz/sidecar

# athenz information
URL=https://localhost:4443/
ATHENZCONF=/home/athenz/bin/linux

# Go related commands
#GOCMD=go
#GOBUILD=$(GOCMD) build
#GOCLEAN=$(GOCMD) clean
#GOTEST=$(GOCMD) test ./...
#GOSYNC=$(GOCMD)  mod tidy


# Redirect error output to a file, so we can show it in development mode.
STDERR=/tmp/.$(PROJECTNAME)-stderr.txt

# PID file will store the server process id when it's running on development mode
PID=/tmp/.$(PROJECTNAME)-api-server.pid

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

# we need to make sure we have go 1.11+
# the output for the go version command is:
# go version go1.11.1 darwin/amd64
GO_VER_GTEQ11 := $(shell expr `go version | cut -f 3 -d' ' | cut -f2 -d.` \>= 12)
ifneq "$(GO_VER_GTEQ11)" "1"
all:
	@echo "Please install 1.12.x or newer version of golang"
endif

# Check richgo does exist.
ifeq (, $(shell which richgo))
$(warning "could not find richgo in $(PATH), run: go get github.com/kyoh86/richgo")
endif

fmt:
	$(info ____________________checking formatting____________________)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

test: sync
	$(info _______________________running tests_______________________)
	 richgo test -v ./...

build: sync
	$(info ________________________building app_______________________)
	GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o agent cmd/server/main.go

sync: 
	$(info _________________downloading dependencies___________________)
	go mod download
	
