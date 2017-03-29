ARCH=$(shell uname -m)
NAME=homehub-cli
VERSION=$(shell cat version.txt)

build:
	rm -rf build
	go build -o build/$(NAME) $(NAME).go

test:
	go test -v github.com/jamesnetherton/homehub-cli/cmd \
	           github.com/jamesnetherton/homehub-cli/cli

release: build
	rm -rf release && mkdir release
	mkdir -p build/linux  && GOOS=linux  go build -ldflags "-X main.version=$(VERSION)" -o build/linux/$(NAME)
	mkdir -p build/darwin && GOOS=darwin go build -ldflags "-X main.version=$(VERSION)" -o build/darwin/$(NAME)
	mkdir -p build/windows && GOOS=windows go build -ldflags "-X main.version=$(VERSION)" -o build/windows/$(NAME).exe

	tar -zcf release/$(NAME)-$(VERSION)-linux-$(ARCH).tgz -C build/linux $(NAME)
	tar -zcf release/$(NAME)-$(VERSION)-darwin-$(ARCH).tgz -C build/darwin $(NAME)
	zip -j release/$(NAME)-$(VERSION)-windows-$(ARCH).zip build/windows/$(NAME).exe

	go get -u github.com/progrium/gh-release
	gh-release checksums sha256
	gh-release create jamesnetherton/$(NAME) $(VERSION)

.PHONY: release build
