# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: blockspiderd transmuted all test clean

GOBIN = ./build/bin
GO ?= latest
GORUN = go run

blockspiderd:
	$(GORUN) build/ci.go install ./cmd/blockspiderd
	@echo "Done building."
	@echo "Run \"$(GOBIN)/blockspiderd\" to launch the blockspider daemon."

transmuted:
	$(GORUN) build/ci.go install ./cmd/transmuted
	@echo "Done building."
	@echo "Run \"$(GOBIN)/transmuted\" to launch the transmute daemon."

all:
	$(GORUN) build/ci.go install

test: all
	$(GORUN) build/ci.go test

lint:
	$(GORUN) build/ci.go lint

clean:
	go clean -cache
	rm -fr build/_workspace/pkg/ $(GOBIN)/*