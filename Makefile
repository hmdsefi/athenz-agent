# project information
PROJECTNAME=$(shell basename "$(PWD)")
GOPKGNAME = gitlab.com/trialblaze/athenz-agent
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
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test ./...
GOSYNC=$(GOCMD)  mod tidy


# Redirect error output to a file, so we can show it in development mode.
STDERR=/tmp/.$(PROJECTNAME)-stderr.txt

# PID file will store the server process id when it's running on development mode
PID=/tmp/.$(PROJECTNAME)-api-server.pid

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

# we need to make sure we have go 1.11+
# the output for the go version command is:
# go version go1.11.1 darwin/amd64
GO_VER_GTEQ11 := $(shell expr `go version | cut -f 3 -d' ' | cut -f2 -d.` \>= 11)
ifneq "$(GO_VER_GTEQ11)" "1"
all:
	@echo "Please install 1.11.x or newer version of golang"
endif

start:
	@echo "  >  $(PROJECTNAME) is available at $(ADDR)"
	@-$(BUILDPATH)/$(PROJECTNAME) 2>&1 & echo $$! > $(PID)
	@cat $(PID) | sed "/^/s/^/  \>  PID: /"

stop:
	@-touch $(PID)
	@-kill `cat $(PID)` 2> /dev/null || true
	@-rm $(PID)

restart: stop start

build:
	@echo "start building $(PROJECTNAME) to $(BUILDPATH)..."
	@echo "create build path..."
	mkdir -p $(BUILDPATH)
	@echo "build path created."
	echo "$(GOBUILD) $(GOBASE)/$(SRC)"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) $(GOBUILD) $(SRC)
	mv $(PROJECTNAME) $(BUILDPATH)
	mkdir $(BUILDPATH)/config
	cp $(GOBASE)/resource/zpe.conf $(BUILDPATH)/config
	cp $(GOBASE)/resource/zpu.conf $(BUILDPATH)/config
	cp $(GOBASE)/agent.json $(BUILDPATH)
	cd $(ATHENZCONF) && ./athenz-conf -o $(BUILDPATH)/config/athenz.conf -z $(URL)
	@echo "build finished."

sync: 
	@echo "checking if there is any missing dependencies..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) $(GOSYNC)
	@echo "synced successfully."
	
