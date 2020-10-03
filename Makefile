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

release-docker:
	docker build -t jamesnetherton/homehub-cli .
	docker tag jamesnetherton/homehub-cli:latest jamesnetherton/homehub-cli:$(VERSION)

	docker login -u "$(DOCKER_USERNAME)" -p "$(DOCKER_PASSWORD)"
	docker push jamesnetherton/homehub-cli:latest
	docker push jamesnetherton/homehub-cli:$(VERSION)
	docker logout

release:
	rm -rf build
	rm -rf release && mkdir release

	mkdir -p build/linux  && GOOS=linux  go build $(BUILDFLAGS) -o build/linux/$(NAME)
	mkdir -p build/rpi32 && GOOS=linux GOARCH=arm go build $(BUILDFLAGS) -o build/rpi32/$(NAME)
	mkdir -p build/rpi64 && GOOS=linux GOARCH=arm64 go build $(BUILDFLAGS) -o build/rpi64/$(NAME)
	mkdir -p build/darwin && GOOS=darwin go build $(BUILDFLAGS) -o build/darwin/$(NAME)
	mkdir -p build/windows && GOOS=windows go build $(BUILDFLAGS) -o build/windows/$(NAME).exe

	tar -zcf release/$(NAME)-$(VERSION)-linux-x86_64.tar.gz -C build/linux $(NAME)
	tar -zcf release/$(NAME)-$(VERSION)-linux-arm.tar.gz -C build/rpi32 $(NAME)
	tar -zcf release/$(NAME)-$(VERSION)-linux-arm64.tar.gz -C build/rpi64 $(NAME)
	tar -zcf release/$(NAME)-$(VERSION)-darwin-x86_64.tar.gz -C build/darwin $(NAME)
	zip -j release/$(NAME)-$(VERSION)-windows-x86_64.zip build/windows/$(NAME).exe

	sha256sum release/$(NAME)-$(VERSION)-linux-x86_64.tar.gz | cut -f1 -d' ' > release/$(NAME)-$(VERSION)-linux-x86_64.tar.gz.sha256
	sha256sum release/$(NAME)-$(VERSION)-linux-arm.tar.gz | cut -f1 -d' ' > release/$(NAME)-$(VERSION)-linux-arm.tar.gz.sha256
	sha256sum release/$(NAME)-$(VERSION)-linux-arm64.tar.gz | cut -f1 -d' ' > release/$(NAME)-$(VERSION)-linux-arm64.tar.gz.sha256
	sha256sum release/$(NAME)-$(VERSION)-darwin-x86_64.tar.gz | cut -f1 -d' ' > release/$(NAME)-$(VERSION)-darwin-x86_64.tar.gz.sha256
	sha256sum release/$(NAME)-$(VERSION)-windows-x86_64.zip | cut -f1 -d' ' > release/$(NAME)-$(VERSION)-windows-x86_64.zip.sha256

.PHONY: release build
