package main

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Local  LocalConfig  `json:"local"`
	Remote RemoteConfig `json:"remote"`
}

type LocalConfig struct {
	SocketPath      string   `json:"socket_path"`
	SshPath         string   `json:"ssh_path"`
	OpPath          string   `json:"op_path"`
	AllowedItemURLs []string `json:"allowed_item_urls"`
}

type RemoteConfig struct {
	SocketPath string `json:"socket_path"`
}

func readConfigFile(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return readConfig(file)
}

func readConfig(r io.Reader) (*Config, error) {
	dec := json.NewDecoder(r)
	var cfg Config
	err := dec.Decode(&cfg)
	return &cfg, err
}
