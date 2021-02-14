VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
            echo v0)
PKGS     = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))


BIN = $(CURDIR)/bin
$(BIN):
	@mkdir -p $@
build: fmt | $(BIN)
	go build \
	    -tags release \
	    -ldflags '-X dasharez0ne-compendium/cmd.Version=$(VERSION)' \
	    -o $(BIN)/dasharez0ne-compendium *.go

run: build
	./bin/dasharez0ne-compendium	

.PHONY: fmt
fmt:
	go fmt .

dasharez0ne.tweets.json:
	twitter-archive -noat -nort @dasharez0ne > $@

update: dasharez0ne.tweets.json
	twitter-archive -noat -nort -a dasharez0ne.tweets.json
	jq . dasharez0ne.tweets.json | jq . | grep media_url\" | sort | cut -f 4 -d '"' | uniq > archive/media_urls.txt
	cd archive && wget --mirror -i media_urls.txt
