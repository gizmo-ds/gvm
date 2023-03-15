NAME=gvm
OUTDIR=build
VERSION=$(shell git describe --tags --always --dirty)
FLAGS=-trimpath -ldflags "-s -w -X main.AppVersion=$(VERSION)"
MAIN=cmd/gvm/main.go
export CGO_ENABLED=0

all: windows-amd64 linux-amd64 darwin-amd64

darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build $(FLAGS) -o $(OUTDIR)/$(NAME)-$@ $(MAIN)

linux-amd64:
	GOOS=linux GOARCH=amd64 go build $(FLAGS) -o $(OUTDIR)/$(NAME)-$@ $(MAIN)

windows-amd64:
	GOOS=windows GOARCH=amd64 go build $(FLAGS) -o $(OUTDIR)/$(NAME)-$@.exe $(MAIN)

sha256sum:
	cd $(OUTDIR); for file in *; do sha256sum $$file > $$file.sha256; done