
proto:
	protoc -I $$GOPATH/src/ -I . otpb/otpb.proto --gofast_out=plugins=grpc:$$GOPATH/src

test:
	go test -v ./...

# build golang bin
build-darwin:
	@GOOS=darwin CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o ./node/bin/opentable-darwin ./main.go
	@chmod +x ./node/bin/opentable-darwin
	@echo "\n\tBinary built for MacOS at ./node/bin/opentable-darwin\n"

build-linux:
	@GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o ./node/bin/opentable-linux ./main.go
	@chmod +x ./node/bin/opentable-linux
	@echo "\n\tBinary built for Linux at ./node/bin/opentable-linux\n"

build: build-linux build-darwin

.PHONY: build test build-linux build-darwin proto
