VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
            echo v0)
PKGS     = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))


BIN = $(CURDIR)/bin
$(BIN):
	@mkdir -p $@
build: fmt | $(BIN) swagger
	go build \
	    -tags release \
	    -ldflags '-X dasharez0ne-compendium/cmd.Version=$(VERSION)' \
	    -o $(BIN)/dasharez0ne-compendium *.go

run: build
	./bin/dasharez0ne-compendium	

.PHONY: fmt
fmt:
	go fmt .

test:
	go test ./...

swagger:
	swag init

deploy:
	fly deploy
