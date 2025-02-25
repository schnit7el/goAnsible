package deployer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/schnit7el/goAnsible/internal/config"
)

type TaskHandler interface {
	Validate(task config.Task) error
	Execute(client *SSHClient, task config.Task) error
}

var taskHandlers = map[string]TaskHandler{
	"command":        CommandHandler{},
	"ssh_setup":      SSHSetupHandler{},
	"file_transfer":  FileTransferHandler{},
	"docker_compose": DockerComposeHandler{},
}

type CommandHandler struct{}

func (h CommandHandler) Validate(task config.Task) error {
	if _, ok := task.Parameters["cmd"]; !ok {
		return fmt.Errorf("missing command")
	}
	return nil
}

func (h CommandHandler) Execute(client *SSHClient, task config.Task) error {
	_, err := client.ExecuteCommand(task.Parameters["cmd"].(string))
	return err
}

type SSHSetupHandler struct{}

func (h SSHSetupHandler) Validate(task config.Task) error {
	required := []string{"public_key_path", "authorized_keys_path"}
	for _, f := range required {
		if _, ok := task.Parameters[f]; !ok {
			return fmt.Errorf("missing %s", f)
		}
	}
	return nil
}

func (h SSHSetupHandler) Execute(client *SSHClient, task config.Task) error {
	pubKey, err := os.ReadFile(task.Parameters["public_key_path"].(string))
	if err != nil {
		return fmt.Errorf("read key failed: %w", err)
	}

	cmd := fmt.Sprintf(
		"mkdir -p $(dirname %s) && echo '%s' >> %s",
		task.Parameters["authorized_keys_path"],
		string(pubKey),
		task.Parameters["authorized_keys_path"],
	)

	_, err = client.ExecuteCommand(cmd)
	return err
}

type FileTransferHandler struct{}

func (h FileTransferHandler) Validate(task config.Task) error {
	required := []string{"src", "dest"}
	for _, f := range required {
		if _, ok := task.Parameters[f]; !ok {
			return fmt.Errorf("missing %s", f)
		}
	}
	return nil
}

func (h FileTransferHandler) Execute(client *SSHClient, task config.Task) error {
	data, err := os.ReadFile(task.Parameters["src"].(string))
	if err != nil {
		return fmt.Errorf("read file failed: %w", err)
	}

	dest := task.Parameters["dest"].(string)
	dirCmd := fmt.Sprintf("mkdir -p %s", filepath.Dir(dest))
	if _, err := client.ExecuteCommand(dirCmd); err != nil {
		return fmt.Errorf("create dir failed: %w", err)
	}

	writeCmd := fmt.Sprintf("echo '%s' | tee %s", string(data), dest)
	if _, err := client.ExecuteCommand(writeCmd); err != nil {
		return fmt.Errorf("write file failed: %w", err)
	}

	return nil
}

type DockerComposeHandler struct{}

func (h DockerComposeHandler) Validate(task config.Task) error {
	required := []string{"action", "file"}
	for _, f := range required {
		if _, ok := task.Parameters[f]; !ok {
			return fmt.Errorf("missing %s", f)
		}
	}
	return nil
}

func (h DockerComposeHandler) Execute(client *SSHClient, task config.Task) error {
	cmd := fmt.Sprintf("docker compose -f %s %s",
		task.Parameters["file"],
		task.Parameters["action"],
	)

	if opts, ok := task.Parameters["options"].([]interface{}); ok {
		for _, opt := range opts {
			cmd += " " + opt.(string)
		}
	}

	_, err := client.ExecuteCommand(cmd)
	return err
}
