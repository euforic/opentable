
# build golang bin
build-mac:
	GOOS=darwin CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o node/scraper ./
	chmod +x node/scraper

build-linux:
	GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o node/scraper ./
	chmod +x node/scraper

.PHONY: build
