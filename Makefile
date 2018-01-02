install:
	go install -v

fmt:
	go fmt
	cd ./cmd && go fmt
	cd ./core && go fmt

certs:
	openssl genrsa \
		-out ./certs/localhost.key \
		2048
	openssl req \
		-new -x509 \
		-key ./certs/localhost.key \
		-out ./certs/localhost.cert \
		-days 3650 \
		-subj /CN=localhost

grpc:
	protoc \
		./messaging/service.proto \
		--gogofaster_out=plugins=grpc:.

.PHONY: fmt install grpc certs
