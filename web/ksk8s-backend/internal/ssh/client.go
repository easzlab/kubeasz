package ssh

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"
)

// Client wraps a crypto/ssh connection.
type Client struct {
	conn   *ssh.Client
	session *ssh.Session
	stdin   io.WriteCloser
	stdout  io.Reader
	stderr  io.Reader
}

// Connect establishes an SSH connection to the local machine (127.0.0.1:22)
// using the default private key or password.
func Connect(addr string, user string, password string, keyPath string) (*Client, error) {
	var authMethods []ssh.AuthMethod

	if password != "" {
		authMethods = append(authMethods, ssh.Password(password))
	}

	if keyPath != "" {
		key, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, fmt.Errorf("read key: %w", err)
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("parse key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("ssh dial: %w", err)
	}

	session, err := conn.NewSession()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("new session: %w", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		session.Close()
		conn.Close()
		return nil, fmt.Errorf("stdin pipe: %w", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		conn.Close()
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		session.Close()
		conn.Close()
		return nil, fmt.Errorf("stderr pipe: %w", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", 80, 24, modes); err != nil {
		session.Close()
		conn.Close()
		return nil, fmt.Errorf("request pty: %w", err)
	}

	if err := session.Shell(); err != nil {
		session.Close()
		conn.Close()
		return nil, fmt.Errorf("start shell: %w", err)
	}

	return &Client{
		conn:    conn,
		session: session,
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
	}, nil
}

// Write sends data to the SSH session's stdin.
func (c *Client) Write(p []byte) (n int, err error) {
	return c.stdin.Write(p)
}

// Read reads from the SSH session's stdout.
func (c *Client) Read(p []byte) (n int, err error) {
	return c.stdout.Read(p)
}

// Resize updates the terminal size.
func (c *Client) Resize(cols, rows int) error {
	return c.session.WindowChange(rows, cols)
}

// Close closes the SSH session and connection.
func (c *Client) Close() error {
	c.session.Close()
	c.conn.Close()
	return nil
}
