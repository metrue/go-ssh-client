# go-ssh-client

This is a little pacakge helps you run command on remote host via SSH

```go
package main

import (
	"log"
	"os"

	ssh "github.com/metrue/go-ssh-client"
)

func main() {
	host := "127.0.0.1"
	script := `
x=1
while [ $x -le 5 ]; do
	echo 'hello'
	x=$(( $x + 1 ))
	sleep 1
done
`
	err := ssh.New(host).
		WithUser("root").
		WithPassword("THEPASSWORDYOUCREATED").
		WithPort("2222").
		RunCommand(script, ssh.CommandOptions{
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			Stdin:  os.Stdin,
		})
	if err != nil {
		log.Fatal(err)
	}
}
```

## Test

```
$ make start_ssh_server
$ make test
$ make clean #clean up running Docker container
```
