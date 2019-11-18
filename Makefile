test:
	go test -v ./...
lint:
	golangci-lint run ./
start_ssh_server:
	@echo "start ssh server ..."
	docker build -t ssh-server -f test/Dockerfile .
	docker run -d --rm --name ssh-server -p 2222:22 ssh-server:latest
clean:
	@echo "stop ssh server ..."
	docker stop ssh-server

all: start_ssh_server test clean
