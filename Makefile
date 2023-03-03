GO_BUILD = go build
GOFLAGS  = CGO_ENABLED=0

## build: Build app binary
.PHONY: build
build:
	$(GOFLAGS) $(GO_BUILD) -a -v -ldflags="-w -s" -o bin/app main.go
