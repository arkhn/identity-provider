include .env example-client/.env
export

.PHONY: build-provider build-client

build-provider:
	go build -o bin/provider

build-client:
	cd example-client && go build -o ../bin/client

run-provider: build-provider
	bin/provider

run-client: build-client
	bin/client