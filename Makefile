.PHONY: proto
proto:
	protoc --proto_path=proto --go_out=plugins=grpc:proto proto/*.proto


.PHONY: build
build:
	go build
