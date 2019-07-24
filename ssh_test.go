package ssh

import (
	"testing"
)

func TestSSH(t *testing.T) {
	host := "127.0.0.1"
	output, err := New(host).
		WithUser("root").
		WithPassword("THEPASSWORDYOUCREATED").
		RunCommand("ls -a")
	if err != nil {
		t.Fatal(err)
	}
	expect := `.
..
.bashrc
.cache
.profile`
	if output != expect {
		t.Fatalf("should get %s but got %s", expect, output)
	}
}
