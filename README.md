# go-ssh-client

This is a little pacakge helps you run command on remote host via SSH

![CI](https://github.com/metrue/go-ssh-client/workflows/ci/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/metrue/go-ssh-client)](https://goreportcard.com/report/github.com/metrue/go-ssh-client)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/metrue/go-ssh-client)
[![asciicast](https://asciinema.org/a/WYvZVCSiAu6FuUksQuhTITIOU.svg)](https://asciinema.org/a/WYvZVCSiAu6FuUksQuhTITIOU)


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
    WithKey("/your/path/to/id_ras").  // Default is ~/.ssh/id_rsa
    WithPort("2222").    // Default is 22
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
