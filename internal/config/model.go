package config

import "time"

type Config struct {
	Nodes []Node `yaml:"nodes"`
	Tasks []Task `yaml:"tasks"`
}

type Node struct {
	Address string `yaml:"address"`
	User    string `yaml:"user"`
	Auth    Auth   `yaml:"auth"`
}

type Auth struct {
	Type             string `yaml:"type"`
	Password         string `yaml:"password,omitempty"`
	SSHKeyPath       string `yaml:"ssh_key_path,omitempty"`
	SSHKeyPassphrase string `yaml:"ssh_key_passphrase,omitempty"`
}

type Task struct {
	Name        string                 `yaml:"name"`
	Type        string                 `yaml:"type"`
	Parameters  map[string]interface{} `yaml:"parameters"`
	Retries     int                    `yaml:"retries,omitempty"`
	Delay       time.Duration          `yaml:"delay,omitempty"`
	IgnoreError bool                   `yaml:"ignore_errors,omitempty"`
}
