VERSION := $(shell git describe --tags)

build:
	GOOS=linux CGO_ENABLED=0  go build -o kubesnoop -ldflags "-X main.version ${VERSION}" server.go

install: build
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ./kubesnoop ${DESTDIR}/usr/local/bin/kubesnoop

test:
	go test ./...

clean:
	rm -f ./kubesnoop

.PHONY: build test install clean
