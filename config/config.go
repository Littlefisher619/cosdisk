package config

import (
	"fmt"
	"github.com/Littlefisher619/cosdisk/repository"
	"github.com/pelletier/go-toml"
)

type CosdiskConfig struct {
	CosURL    string
	SecretID  string
	SecretKey string
}

type Config struct {
	repository.RepositoryConfig
	CosdiskConfig
}

func loadCosdiskConfig(config *toml.Tree) (CosdiskConfig, error) {
	cosURL := config.Get("storage.cosURL")
	if cosURL == nil {
		return CosdiskConfig{}, fmt.Errorf("storage.cosURL not found")
	}
	secretID := config.GetDefault("storage.secretID", "")
	secretKey := config.GetDefault("storage.secretKey", "")
	return CosdiskConfig{
		CosURL:    cosURL.(string),
		SecretID:  secretID.(string),
		SecretKey: secretKey.(string),
	}, nil
}

func LoadConfig(configPath string) (*Config, error) {
	config, err := toml.LoadFile(configPath)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	c.RepositoryConfig, err = repository.LoadDataBaseConfig(config)
	if err != nil {
		return nil, err
	}
	c.CosdiskConfig, err = loadCosdiskConfig(config)
	if err != nil {
		return nil, err
	}
	return c, nil
}
