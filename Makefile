VERSION = $(shell date -u +.%Y%m%d.%H%M%S)

all: export GOPATH=${PWD}/Godeps/_workspace:${PWD}/../../../..
all: format
	@mkdir -p bin/
	@echo "--> Running go build"
	@go build -ldflags "-X github.com/untoldwind/gorrd/config.versionMinor=${VERSION}" -v -o bin/gorrd github.com/untoldwind/gorrd

format: export GOPATH=${PWD}/Godeps/_workspace:${PWD}/../../../..
format:
	@echo "--> Running go fmt"
	@go fmt ./...

godepssave:
	@echo "--> Godeps save"
	@go build -v -o bin/godep github.com/tools/godep
	@bin/godep save

genctags:
	@echo "--> Gen Ctags"
	@go build -v -o bin/gotags github.com/jstemmer/gotags
	@bin/gotags -tag-relative=true -R=true -sort=true -f=".tags" -fields=+l .
