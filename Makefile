
# build golang bin
build-darwin:
	@GOOS=darwin CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o ./bin/opentable-darwin ./opentable/cmd/.
	@chmod +x bin/opentable-darwin
	@echo "\n\tBinary built for MacOS at ./bin/opentable-darwin\n"

build-linux:
	@GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o ./bin/opentable-linux ./opentable/cmd/.
	@chmod +x bin/opentable-linux
	@echo "\n\tBinary built for Linux at ./bin/opentable-linux\n"

build: build-linux build-darwin

test:
	@echo "No tests available"

.PHONY: build test build-linux build-darwin
