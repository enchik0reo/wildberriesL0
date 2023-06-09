package config

import (
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Port      string `yaml:"port" env-default:"8080"`
	CacheSize int    `yaml:"cache_size"`
	Nats      struct {
		URL       string `yaml:"url"`
		ClusterID string `yaml:"clusterid"`
		ClientID  string `yaml:"clientid"`
	} `yaml:"nats"`
	DB struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `env:"POSTGRES_PASSWORD"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"db"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	p := filepath.Join("config", "cfg.yml")
	if err := cleanenv.ReadConfig(p, cfg); err != nil {
		_, err = cleanenv.GetDescription(cfg, nil)
		return nil, err
	}

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg.DB.Password = os.Getenv("POSTGRES_PASSWORD")

	return cfg, nil
}
