SRC=$(shell find . -name "*.go")

# project information
PROJECTNAME=$(shell basename "$(PWD)")
PKG_DATE=$(shell date '+%Y-%m-%dT%H:%M:%S')
SRC=cmd/tools/athenz-agent.go
GOBASE=$(PWD)
GOBIN=$(GOBASE)/bin


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

codecov: sync
	$(info __________________running tests coverage___________________)
	 sh build/script/coverage.sh

build: sync
	$(info ________________________building app_______________________)
	GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o agent cmd/server/main.go

sync: 
	$(info _________________downloading dependencies___________________)
	go get -v ./...
