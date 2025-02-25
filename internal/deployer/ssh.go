package deployer

import (
	"fmt"
	"os"
	"time"

	"github.com/schnit7el/goAnsible/internal/config"
	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	*ssh.Client
}

func NewSSHClient(node config.Node) (*SSHClient, error) {
	cfg := &ssh.ClientConfig{
		User:            node.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	if node.Auth.SSHKeyPath != "" {
		keyAuth, err := parseSSHKey(
			node.Auth.SSHKeyPath,
			[]byte(node.Auth.SSHKeyPassphrase),
		)
		if err != nil {
			return nil, fmt.Errorf("ssh key auth failed: %w", err)
		}
		cfg.Auth = append(cfg.Auth, keyAuth)
	}

	if node.Auth.Password != "" {
		cfg.Auth = append(cfg.Auth, ssh.Password(node.Auth.Password))
	}

	conn, err := ssh.Dial("tcp", node.Address, cfg)
	if err != nil {
		return nil, fmt.Errorf("ssh connection failed: %w", err)
	}

	return &SSHClient{conn}, nil
}

func (c *SSHClient) ExecuteCommand(cmd string) (string, error) {
	session, err := c.Client.NewSession()
	if err != nil {
		return "", fmt.Errorf("session failed: %w", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	return string(output), err
}

func parseSSHKey(path string, passphrase []byte) (ssh.AuthMethod, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key: %w", err)
	}

	key, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		if _, ok := err.(*ssh.PassphraseMissingError); ok {
			key, err = ssh.ParsePrivateKeyWithPassphrase(keyBytes, passphrase)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("key parse failed: %w", err)
	}

	return ssh.PublicKeys(key), nil
}
