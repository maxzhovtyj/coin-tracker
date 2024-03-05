package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

type Config struct {
	TelegramApiToken string  `yaml:"telegramApiToken"`
	DBPath           string  `yaml:"dbPath"`
	Binance          Binance `yaml:"binance"`
}

type Binance struct {
	ApiKey    string `yaml:"apiKey"`
	SecretKey string `yaml:"secretKey"`
}

var (
	configPath = flag.String("config", "", "config path")

	cfg  *Config
	once sync.Once
)

func New() (*Config, error) {
	var err error

	once.Do(func() {
		flag.Parse()

		if *configPath == "" {
			err = fmt.Errorf("empty config path")
			return
		}

		var rawConfig []byte
		rawConfig, err = os.ReadFile(*configPath)
		if err != nil {
			err = fmt.Errorf("can't read config file: %w", err)
			return
		}

		err = yaml.Unmarshal(rawConfig, &cfg)
		if err != nil {
			err = fmt.Errorf("failed to unmarshal config file: %w", err)
			return
		}
	})
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
