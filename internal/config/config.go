package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

type Config struct {
	TelegramApiToken string `yaml:"telegramApiToken"`
}

var (
	configPath = flag.String("config", "", "config path")

	cfg  *Config
	once sync.Once
)

func New() (*Config, error) {
	flag.Parse()

	if *configPath == "" {
		return nil, fmt.Errorf("empty config path")
	}

	rawConfig, err := os.ReadFile(*configPath)
	if err != nil {
		return nil, fmt.Errorf("can't read config file: %w", err)
	}

	var parseErr error
	once.Do(func() {
		parseErr = yaml.Unmarshal(rawConfig, &cfg)
	})
	if parseErr != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", parseErr)
	}

	return cfg, nil
}
