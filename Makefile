VERSION := $(shell git describe --tags)

build:
	GOOS=linux CGO_ENABLED=0  go build -o kubesnoop -a -installsuffix cgo server.go

install: build
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ./kubesnoop ${DESTDIR}/usr/local/bin/kubesnoop

test:
	go test ./...

clean:
	rm -f ./kubesnoop

deploy: build
	git commit -m "Latest build" kubesnoop
	git push deis master

.PHONY: build test install clean
