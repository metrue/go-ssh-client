package ssh

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// Client ssh client
type Client struct {
	server string
	port   string
	user   string

	key      string
	password string

	session *ssh.Session
	conn    ssh.Conn
}

// New create a client
func New(server string) Client {
	return Client{
		server: server,
		port:   "22",
	}
}

// WithServer with server
func (c Client) WithServer(addr string) Client {
	return Client{
		server:   addr,
		port:     c.port,
		user:     c.user,
		key:      c.key,
		password: c.password,
	}
}

// WithUser with key
func (c Client) WithUser(user string) Client {
	return Client{
		server:   c.server,
		port:     c.port,
		user:     user,
		key:      c.key,
		password: c.password,
	}
}

// WithPassword with key
func (c Client) WithPassword(password string) Client {
	return Client{
		server:   c.server,
		port:     c.port,
		user:     c.user,
		key:      c.key,
		password: password,
	}
}

// WithKey with key
func (c Client) WithKey(keyfile string) Client {
	return Client{
		server:   c.server,
		port:     c.port,
		user:     c.user,
		key:      keyfile,
		password: c.password,
	}
}

// WithPort with port
func (c Client) WithPort(port string) Client {
	return Client{
		server:   c.server,
		port:     port,
		user:     c.user,
		key:      c.key,
		password: c.password,
	}
}

// RunCommand run command onto remote server via SSH
func (c Client) RunCommand(command string) ([]byte, []byte, error) {
	client, err := c.connect()
	if err != nil {
		return nil, nil, err
	}

	defer func() {
		if err := client.disconnect(); err != nil {
			fmt.Println("-->", err)
			log.Println(err)
		}
	}()

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	client.session.Stdout = &stdoutBuf
	client.session.Stderr = &stderrBuf

	err = client.session.Run(command)
	return stdoutBuf.Bytes(), stderrBuf.Bytes(), err
}

// Connect connect server
func (c Client) connect() (Client, error) {
	Auth := []ssh.AuthMethod{}
	if c.key != "" {
		publicKey, err := publicKey(c.key)
		if err != nil {
			return Client{}, err
		}
		Auth = append(Auth, publicKey)
	}

	if c.password != "" {
		Auth = append(Auth, ssh.Password(c.password))
	}

	config := &ssh.ClientConfig{
		User:            c.user,
		Auth:            Auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil },
	}

	addr := net.JoinHostPort(c.server, c.port)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return Client{}, err
	}

	session, err := conn.NewSession()
	if err != nil {
		return Client{}, err
	}

	return Client{
		server:   c.server,
		port:     c.port,
		user:     c.user,
		key:      c.key,
		password: c.password,

		conn:    conn,
		session: session,
	}, nil
}

// Disconnect disconnect with server
func (c Client) disconnect() error {
	if err := c.session.Close(); err != nil {
		// "https://github.com/golang/go/issues/28108"
		if err == io.EOF {
			return nil
		}
		return errors.Wrap(err, "session close failure")
	}
	if err := c.conn.Close(); err != nil {
		return errors.Wrap(err, "connection close failure")
	}
	return nil
}

func publicKey(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}
