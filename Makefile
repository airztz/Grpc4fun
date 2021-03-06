.PHONY: hello-server
hello-server:
	go build -mod=vendor -o target/hello-server servers/grpc/main.go
	chmod a+x target/hello-server

.PHONY: hello-client
hello-client:
	go build -mod=vendor -o target/hello-client clients/grpc/main.go
	chmod a+x target/hello-client

.PHONY: rest-hello-server
rest-hello-server:
	go build -mod=vendor -o target/rest-hello-server servers/rest/main.go
	chmod a+x target/rest-hello-server

.PHONY: clean
clean:
	rm -rf target

.PHONY: build
build: hello-server hello-client rest-hello-server