.PHONY: hello-server
hello-server:
	go build -mod=vendor -o target/hello-server servers/hello/main.go
	chmod a+x target/hello-server

.PHONY: hello-client
hello-client:
	go build -mod=vendor -o target/hello-client clients/hello/main.go
	chmod a+x target/hello-client

.PHONY: clean
clean:
	rm -rf target

.PHONY: build
build: hello-server hello-client