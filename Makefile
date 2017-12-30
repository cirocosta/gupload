install:
	go install -v

fmt:
	go fmt
	cd ./cmd && go fmt
	cd ./core && go fmt

gen:
	protoc \
		./messaging/service.proto \
		--gogofaster_out=plugins=grpc:.

.PHONY: gen
