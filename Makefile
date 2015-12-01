SHELL := /bin/bash
PKG = github.com/Clever/riemann-gearman
SUBPKGS =
PKGS = $(PKG) $(SUBPKGS)

.PHONY: test $(PKGS)

test: $(PKG)

$(PKG):
	go get github.com/golang/lint/golint
	$(GOPATH)/bin/golint $(GOPATH)/src/$@*/**.go
	go get -d -t $@
	go test -cover -coverprofile=$(GOPATH)/src/$@/c.out $@ -test.v
ifeq ($(HTMLCOV),1)
	go tool cover -html=$(GOPATH)/src/$@/c.out
endif


SHELL := /bin/bash
PKGS := $(shell go list ./... | grep -v /vendor)
GODEP := $(GOPATH)/bin/godep

$(GODEP):
	go get -u github.com/tools/godep

vendor: $(GODEP)
	$(GODEP) save $(PKGS)
	find vendor/ -path '*/vendor' -type d | xargs -IX rm -r X # remove any nested vendor directories
