start_ssh_server:
	docker build -t ssh-server -f test/Dockerfile .
	docker run -d --rm --name ssh-server -p 22:22 ssh-server:latest
test:
	go test -v ./...
lint:
	golangci-lint run ./
clean:
	docker stop ssh-server
