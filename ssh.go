package ssh

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// Clienter defines interface of SSH client
type Clienter interface {
	WithServer(add string) Client
	WithUser(user string) Client
	WithPassword(password string) Client
	WithKey(key string) Client
	WithPort(port string) Client
	Connectable(timeout time.Duration) (bool, error)
	RunCommand(command string, options CommandOptions) error
}

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
	home, _ := homedir.Dir()
	return Client{
		server: server,
		user:   "root",
		port:   "22",
		key:    filepath.Join(home, ".ssh/id_rsa"),
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

// CommandOptions options for command
type CommandOptions struct {
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	Timeout time.Duration
	Env     []string
}

// RunCommand run command onto remote server via SSH
func (c Client) RunCommand(command string, options CommandOptions) error {
	timeout := 20 * time.Second
	if options.Timeout > 0 {
		timeout = options.Timeout
	}
	client, err := c.connect(timeout)
	if err != nil {
		return err
	}

	defer func() {
		if err := client.disconnect(); err != nil {
			fmt.Println("-->", err)
			log.Println(err)
		}
	}()

	if options.Stdin != nil {
		stdin, err := client.session.StdinPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdin for session: %v", err)
		}
		// nolint
		go io.Copy(stdin, options.Stdin)
	}

	if options.Stdout != nil {
		stdout, err := client.session.StdoutPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdout for session: %v", err)
		}
		// nolint
		go io.Copy(options.Stdout, stdout)
	}

	if options.Stderr != nil {
		stderr, err := client.session.StderrPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stderr for session: %v", err)
		}
		// nolint
		go io.Copy(options.Stderr, stderr)
	}

	for _, env := range options.Env {
		variable := strings.Split(env, "=")
		if len(variable) != 2 {
			continue
		}

		if err := client.session.Setenv(variable[0], variable[1]); err != nil {
			return err
		}
	}

	return client.session.Run(command)
}

// Connectable check if client can connect to ssh server within timeout
func (c Client) Connectable(timeout time.Duration) (bool, error) {
	client, err := c.connect(timeout)
	if err != nil {
		return false, err
	}

	defer func() {
		if err := client.disconnect(); err != nil {
			fmt.Println("-->", err)
			log.Println(err)
		}
	}()

	return true, nil
}

// Connect connect server
func (c Client) connect(timeout time.Duration) (Client, error) {
	Auth := []ssh.AuthMethod{}

	if c.password != "" {
		Auth = append(Auth, ssh.Password(c.password))
	} else if c.key != "" {
		publicKey, err := publicKey(c.key)
		if err != nil {
			return Client{}, err
		}
		Auth = append(Auth, publicKey)
	} else {
		return Client{}, fmt.Errorf("password or keyfile required for ssh connection ")
	}

	config := &ssh.ClientConfig{
		User:            c.user,
		Auth:            Auth,
		Timeout:         timeout,
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
	if err != nil && err.Error() == "ssh: this private key is passphrase protected" {
		// Note: Newer version of Go support the error type `PassphraseMissingError'. Using that, the error check is:
		// if err, ok := err.(*PassphraseMissingError); !ok {...}
		// Note: Put the variable name just into a variable; would need to be passed through from fx/middlewares/ssh.go
		// using context.Contexter.
		envName := "SSH_PASS_PHRASE"
		passphrase := os.Getenv(envName)
		if passphrase != "" {
			key, err = ssh.ParsePrivateKeyWithPassphrase(buffer, []byte(passphrase))
			if err != nil {
				err = fmt.Errorf("Using passphrase from environment: %s, %v", envName, err)
			}
		} else {
			err = fmt.Errorf("No passphrase defined by environment: %s, %v", envName, err)
		}
	}
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

var (
	_ Clienter = Client{}
)
