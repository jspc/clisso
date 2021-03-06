GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=clisso
VERSION=`git describe --tags --always`

.PHONY: build
build:
	$(GOBUILD) -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME) -v

.PHONY: test
test:
	$(GOCMD) test -v ./...

.PHONY: darwin-386
darwin-386:
	GOOS=darwin GOARCH=386 $(GOBUILD) -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME)-darwin-386 -v

.PHONY: darwin-amd64
darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME)-darwin-amd64 -v

.PHONY: linux-386
linux-386:
	GOOS=linux GOARCH=386 $(GOBUILD) -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME)-linux-386 -v

.PHONY: linux-amd64
linux-amd64:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME)-linux-amd64 -v

.PHONY: windows-386
windows-386:
	GOOS=windows GOARCH=386 $(GOBUILD) -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME)-windows-386.exe -v

.PHONY: windows-amd64
windows-amd64:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME)-windows-amd64.exe -v

.PHONY: all
all: darwin-386 darwin-amd64 linux-386 linux-amd64 windows-386 windows-amd64

.PHONY: zip
zip:
	for i in `ls -1 $(BINARY_NAME)* | grep -v '.zip'`; do zip $$i.zip $$i; done

.PHONY: release
release: clean all zip

.PHONY: install
install:
	go install -ldflags "-X main.version=$(VERSION)"

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)*
