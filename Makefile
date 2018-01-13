ARCH=$(shell uname -m)
NAME=homehub-cli
ROOT_PACKAGE := $(shell go list ./cmd)
FIRMWARE=$(shell cat firmware.txt)
VERSION=$(shell cat version.txt)
REVISION=$(shell git rev-parse --short HEAD || echo 'Unknown')
BUILD_DATE=$(shell date +%d/%m/%Y)
BUILDFLAGS := -ldflags \
  " -X $(ROOT_PACKAGE).version=$(VERSION)\
    -X $(ROOT_PACKAGE).firmware=$(FIRMWARE)\
    -X $(ROOT_PACKAGE).revision=$(REVISION)\
    -X $(ROOT_PACKAGE).date=$(BUILD_DATE)"

build:
	rm -rf build
	go build $(BUILDFLAGS) -o build/$(NAME) $(NAME).go

test: build
	go test -v github.com/jamesnetherton/homehub-cli/cmd \
	           github.com/jamesnetherton/homehub-cli/cli

docker:
	docker build -t jamesnetherton/homehub-cli .
	docker tag jamesnetherton/homehub-cli:latest jamesnetherton/homehub-cli:$(VERSION)

release: docker
	rm -rf build
	rm -rf release && mkdir release
	mkdir -p build/linux  && GOOS=linux  go build $(BUILDFLAGS) -o build/linux/$(NAME)
	mkdir -p build/darwin && GOOS=darwin go build $(BUILDFLAGS) -o build/darwin/$(NAME)
	mkdir -p build/windows && GOOS=windows go build $(BUILDFLAGS) -o build/windows/$(NAME).exe

	tar -zcf release/$(NAME)-$(VERSION)-linux-$(ARCH).tar.gz -C build/linux $(NAME)
	tar -zcf release/$(NAME)-$(VERSION)-darwin-$(ARCH).tar.gz -C build/darwin $(NAME)
	zip -j release/$(NAME)-$(VERSION)-windows-$(ARCH).zip build/windows/$(NAME).exe

	go get -u github.com/progrium/gh-release
	gh-release checksums sha256
	gh-release create jamesnetherton/$(NAME) $(VERSION)

.PHONY: release build
