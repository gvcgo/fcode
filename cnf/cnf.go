package cnf

import (
	"fmt"
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

var (
	DefaultConf  *Conf
	DefaultModel *AIModel
)

func init() {
	DefaultConf = &Conf{}
	DefaultConf.load()
	if len(DefaultConf.AIModels) > 0 {
		DefaultModel = DefaultConf.AIModels[0]
	}
}

const (
	FCodeDir        = ".fcode"
	FCodeConfigFile = "conf.toml"
	FCodeApiKeyFile = "key.toml"
	DefaultCursor   = "<CURSOR>"
	DefaultPort     = 8123
	ModelCtxKey     = "ai_model"
)

type AIModel struct {
	Name     string `toml:"name"`
	Type     string `toml:"type"`
	Api      string `toml:"api"`
	Model    string `toml:"model"`
	Key      string `toml:"key"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type Conf struct {
	Cursor   string     `toml:"cursor"`
	Port     int        `toml:"port"`
	AIModels []*AIModel `toml:"models"`
}

func (c *Conf) GetPath() string {
	homeDir, _ := os.UserHomeDir()
	fcodeDir := filepath.Join(homeDir, FCodeDir)
	os.MkdirAll(fcodeDir, os.ModePerm)
	return filepath.Join(fcodeDir, FCodeConfigFile)
}

func (c *Conf) load() {
	path := c.GetPath()
	content, _ := os.ReadFile(path)
	if len(content) > 0 {
		_ = toml.Unmarshal(content, c)
	}
}

func (c *Conf) save() {
	path := c.GetPath()
	content, err := toml.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	os.WriteFile(path, content, os.ModePerm)
	c.load()
}

func (c *Conf) GetCursor() string {
	if c.Cursor != "" {
		return c.Cursor
	}
	return DefaultCursor
}

func (c *Conf) GetPort() string {
	if c.Port > 0 {
		return fmt.Sprintf(":%d", c.Port)
	}
	return fmt.Sprintf(":%d", DefaultPort)
}

func (c *Conf) SetApiKey(name, key string) {
	c.load()
	for _, m := range c.AIModels {
		if m.Name == name {
			m.Key = key
		}
	}
	c.save()
}
