package main

import (
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

var (
	DefaultConf *Conf
	DefaultKey  *Key
)

func init() {
	DefaultConf = &Conf{}
	DefaultConf.load()
	DefaultKey = &Key{}
	DefaultKey.Load()
}

const (
	FCodeDir        = ".fcode"
	FCodeConfigFile = "conf.toml"
	FCodeApiKeyFile = "key.toml"
	IdeName         = "vim"
	PluginVersion   = "0.2.1"
)

type Conf struct {
	UserName string `toml:"username"`
	Password string `toml:"password"`
}

func (c *Conf) getPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, FCodeDir, FCodeConfigFile)
}

func (c *Conf) load() {
	path := c.getPath()
	content, _ := os.ReadFile(path)
	if len(content) > 0 {
		_ = toml.Unmarshal(content, c)
	}
}

type Key struct {
	APIKey string `toml:"key"`
}

func (k *Key) getPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, FCodeDir, FCodeApiKeyFile)
}

func (k *Key) Load() {
	path := k.getPath()
	content, _ := os.ReadFile(path)
	if len(content) > 0 {
		_ = toml.Unmarshal(content, k)
	}
	if k.APIKey == "" {
		Login()
	}
}

func (k *Key) Save() {
	path := k.getPath()
	content, _ := toml.Marshal(k)
	os.WriteFile(path, content, os.ModePerm)
}
