# go-ssh-client

This is a little pacakge helps you run command on remote host via SSH

```
	host := "127.0.0.1"
	output, err := New(host).
		WithUser("root").
		WithPassword("THEPASSWORDYOUCREATED").
		RunCommand("ls -a")
	if err != nil {
		t.Fatal(err)
	}
```
