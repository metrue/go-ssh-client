package ssh

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func TestSSH(t *testing.T) {
	t.Run("public key", func(t *testing.T) {
		cases := []struct {
			cmd           string
			errIsNil      bool
			stdoutIsEmpty bool
			stderrIsEmpty bool
		}{
			{
				cmd:           "ls -a",
				errIsNil:      true,
				stdoutIsEmpty: false,
				stderrIsEmpty: true,
			},
			{
				cmd:           "docker ps",
				errIsNil:      false,
				stdoutIsEmpty: true,
				stderrIsEmpty: false,
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
			err := New(host).
				WithUser("root").
				WithPort("2222").
				WithKey("./test/id_rsa").
				RunCommand(c.cmd, options)

			if !reflect.DeepEqual(err == nil, c.errIsNil) {
				t.Fatalf("should get %v but got %v", c.errIsNil, err == nil)
			}

			if !reflect.DeepEqual(len(outPipe.String()) == 0, c.stdoutIsEmpty) {
				t.Fatalf("should get %v but got %v", c.stdoutIsEmpty, len(outPipe.String()) == 0)
			}

			if !reflect.DeepEqual(len(errPipe.String()) == 0, c.stderrIsEmpty) {
				t.Fatalf("should get %v but got %v", c.stderrIsEmpty, len(errPipe.String()) == 0)
			}
		}
	})

	t.Run("password", func(t *testing.T) {

		cases := []struct {
			cmd           string
			errIsNil      bool
			stdoutIsEmpty bool
			stderrIsEmpty bool
		}{
			{
				cmd:           "ls -a",
				errIsNil:      true,
				stdoutIsEmpty: false,
				stderrIsEmpty: true,
			},
			{
				cmd:           "docker ps",
				errIsNil:      false,
				stdoutIsEmpty: true,
				stderrIsEmpty: false,
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

			err := New(host).
				WithUser("root").
				WithPort("2222").
				WithPassword("THEPASSWORDYOUCREATED").
				RunCommand(c.cmd, options)

			if !reflect.DeepEqual(err == nil, c.errIsNil) {
				t.Fatalf("should get %v but got %v", c.errIsNil, err == nil)
			}

			if !reflect.DeepEqual(len(outPipe.String()) == 0, c.stdoutIsEmpty) {
				t.Fatalf("should get %v but got %v", c.stdoutIsEmpty, len(outPipe.String()) == 0)
			}

			if !reflect.DeepEqual(len(errPipe.String()) == 0, c.stderrIsEmpty) {
				t.Fatalf("should get %v but got %v", c.stderrIsEmpty, len(errPipe.String()) == 0)
			}
		}
	})
}
