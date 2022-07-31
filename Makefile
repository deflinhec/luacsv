# Define
VERSION=1.0.0
BUILD=$(shell git rev-parse HEAD)

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS='-s -w -X "main.Version=$(VERSION)" -X "main.Build=$(BUILD)"'

.PHONY: default all assets

default: all

fmt:
	go fmt ./...

tidy:
	go mod tidy

go-bindata:
	go install github.com/go-bindata/go-bindata/...

assets: go-bindata
	go-bindata -nomemcopy -pkg=assets -o=assets/assets.go \
		-debug=$(if $(findstring debug,$(BUILDTAGS)),true,false) \
		-ignore=assets.go assets/...

# Sperate "linux-amd64" as GOOS and GOARCH
OSARCH_SPERATOR = $(word $2,$(subst -, ,$1))
# Platform build options
cross-compile-%: export GOOS=$(call OSARCH_SPERATOR,$*,1)
cross-compile-%: export GOARCH=$(call OSARCH_SPERATOR,$*,2)
cross-compile-%: assets
	go build -ldflags $(LDFLAGS) -o ./build/$(GOOS)-$(GOARCH)/ ./cmd/...

linux: cross-compile-linux-amd64
darwin: cross-compile-darwin-amd64
windows: cross-compile-windows-amd64

# Local build options
build: assets
	go build -ldflags $(LDFLAGS) ./cmd/...

all: fmt tidy darwin linux windows
