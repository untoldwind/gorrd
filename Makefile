VERSION = $(shell date -u +.%Y%m%d.%H%M%S)

all: export GOPATH=${PWD}/../../../..
all: format
	@mkdir -p bin/
	@echo "--> Running go build"
	@go build -ldflags "-X github.com/untoldwind/gorrd/config.versionMinor=${VERSION}" -v -o bin/gorrd github.com/untoldwind/gorrd

format: export GOPATH=${PWD}/../../../..
format:
	@echo "--> Running go fmt"
	@go fmt ./...

test: export GOPATH=${PWD}/../../../..
test:
	@echo "--> Running tests"
	@go test -v ./...
	@$(MAKE) vet

godepssave:
	@echo "--> Godeps save"
	@go build -v -o bin/godep github.com/tools/godep
	@bin/godep save ./...

genctags:
	@echo "--> Gen Ctags"
	@go build -v -o bin/gotags github.com/jstemmer/gotags
	@bin/gotags -tag-relative=true -R=true -sort=true -f=".tags" -fields=+l .

goconvey:
	@go build -v -o bin/goconvey github.com/smartystreets/goconvey
	@bin/goconvey

vet: export GOPATH=${PWD}/Godeps/_workspace:${PWD}/../../../..
vet:
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
    	go get golang.org/x/tools/cmd/vet; \
    fi
	@echo "--> Running go tool vet $(VETARGS)"
	@find . -name "*.go" | grep -v "./vendor/" | xargs go tool vet $(VETARGS); if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for reviewal."; \
	fi
