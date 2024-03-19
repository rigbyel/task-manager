package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	defaultConfigPath = "./config/local.yaml"
)

type Config struct {
	Env         string `yaml:"env" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server" env-required:"true"`
}

type HTTPServer struct {
	Address      string        `yaml:"address" env-required:"true"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"3s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"3s"`
}

// loading config from configPath
func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file doesn't exist:" + configPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		panic("unable to read config file")
	}

	return &cfg
}

// fetch config path
// priority: command line flags (--config="pathtoconfig") > environmental variables > default
func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG")
	}

	if path == "" {
		path = defaultConfigPath
	}

	return path
}
