package ssh

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestSSH(t *testing.T) {
	t.Run("Connectable", func(t *testing.T) {
		cases := []struct {
			host        string
			connectable bool
		}{{
			host:        "127.0.0.1",
			connectable: true,
		},
			{
				host:        "192.3.0.1",
				connectable: false,
			},
		}
		for _, c := range cases {
			client := New(c.host).
				WithUser("root").
				WithPort("2222").
				WithKey("./test/id_rsa")
			fmt.Println(client)
			ok, err := client.Connectable(5 * time.Second)
			if !ok {
				if err == nil {
					t.Fatalf("error should not be nil ")
				}
			}
			if ok != c.connectable {
				t.Fatalf("should get %v but got %v", c.connectable, ok)
			}
		}
	})

	t.Run("public key", func(t *testing.T) {
		cases := []struct {
			cmd    string
			stdout string
			stderr string
		}{
			{
				cmd:    "echo 1",
				stdout: "1\n",
				stderr: "",
			},
			{
				cmd:    "docker ps",
				stdout: "",
				stderr: "bash: docker: command not found\n",
			},
		}

		for _, c := range cases {
			host := "127.0.0.1"
			var inPipe bytes.Buffer
			var outPipe bytes.Buffer
			var errPipe bytes.Buffer
			options := CommandOptions{
				Stdin:  bufio.NewReader(&inPipe),
				Stdout: bufio.NewWriter(&outPipe),
				Stderr: bufio.NewWriter(&errPipe),
			}
			_ = New(host).
				WithUser("root").
				WithPort("2222").
				WithKey("./test/id_rsa").
				RunCommand(c.cmd, options)

			if errPipe.String() != c.stderr {
				t.Fatalf("should get %v but got %v", c.stderr, errPipe.String())
			}
			if outPipe.String() != c.stdout {
				t.Fatalf("should get %v but got %v", c.stdout, outPipe.String())
			}
		}
	})

	t.Run("password", func(t *testing.T) {
		cases := []struct {
			cmd    string
			stdout string
			stderr string
		}{
			{
				cmd:    "echo 1",
				stdout: "1\n",
				stderr: "",
			},
			{
				cmd:    "docker ps",
				stdout: "",
				stderr: "bash: docker: command not found\n",
			},
		}

		for _, c := range cases {
			host := "127.0.0.1"
			var inPipe bytes.Buffer
			var outPipe bytes.Buffer
			var errPipe bytes.Buffer
			options := CommandOptions{
				Stdin:  bufio.NewReader(&inPipe),
				Stdout: bufio.NewWriter(&outPipe),
				Stderr: bufio.NewWriter(&errPipe),
			}
			_ = New(host).
				WithUser("root").
				WithPort("2222").
				WithPassword("THEPASSWORDYOUCREATED").
				RunCommand(c.cmd, options)

			if errPipe.String() != c.stderr {
				t.Fatalf("should get %v but got %v", c.stderr, errPipe.String())
			}

			if outPipe.String() != c.stdout {
				t.Fatalf("should get %v but got %v", c.stdout, outPipe.String())
			}
		}
	})
}
