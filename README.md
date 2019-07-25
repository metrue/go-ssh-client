# go-ssh-client

This is a little pacakge helps you run command on remote host via SSH

```go
package main

import (
	"log"

	ssh "github.com/metrue/go-ssh-client"
)

func main() {
	host := "127.0.0.1"
	output, err := ssh.New(host).
		WithUser("root").
		WithPassword("THEPASSWORDYOUCREATED").
		RunCommand("ls -a")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(output)
	}
}
```

## Test

```
$ make start_ssh_server
$ make test
$ make clean #clean up running Docker container
```
