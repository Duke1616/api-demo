package conf

import (
	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
)

// LoadConfigFromEnv 从环境变量中加载
func LoadConfigFromEnv() error {
	cfg := NewDefaultConfig()
	if err := env.Parse(cfg); err != nil {
		return err
	}
	return nil
}

// LoadConfigFromToml 从文件中加载
func LoadConfigFromToml(path string) error {
	cfg := NewDefaultConfig()
	_, err := toml.DecodeFile(path, cfg)
	if err != nil {
		return err
	}
	SetGlobalConfig(cfg)
	return nil
}
