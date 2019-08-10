.PHONY: all

all: test clean

test: start_ssh_server
	go test -v ./...
lint:
	golangci-lint run ./
start_ssh_server:
	@echo "start ssh server ..."
	docker build -t ssh-server -f test/Dockerfile .
	docker run -d --rm --name ssh-server -p 22:22 ssh-server:latest
clean:
	@echo "stop ssh server ..."
	docker stop ssh-server
