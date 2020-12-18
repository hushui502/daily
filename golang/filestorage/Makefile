VERSION := v0.0.1

.PONY: build install release release-linux release-mac

run:
	@go build
	@./filestorage

build:
	@go build -ldflags '-w -s' -o filestorage

install:
	@cp -f filestorage /usr/local/bin

release: release-linux release-mac

release-linux:
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o filestorage
	@tar zcvf filestorage-$(VERSION)-linux-amd64.tar.gz ./filestorage
	@rm -f filestorage

release-mac:
	@env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags '-w -s' -o filestorage
	@tar zcvf filestorage-$(VERSION)-darwin-amd64.tar.gz ./filestorage
	@rm -f filestorage