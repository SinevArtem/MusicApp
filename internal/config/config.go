package config

import (
	"flag"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
	Redis       `yaml:"redis"`
}

type HTTPServer struct {
	Address string        `yaml:"address" env-default:"localhost:8080"`
	Timeout time.Duration `yaml:"timeout"`
}

type Redis struct {
	Address  string `yaml:"address" env-default:"localhost:6379"`
	Password string `yaml:"password" env-default:""`
	DB       int    `yaml:"db" env-default:"0"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		log.Fatal("not found config path")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("cannot red config: %s", err)
	}

	return &cfg
}

func fetchConfigPath() string {
	var path string

	// --config="path/to/config.yaml"
	flag.StringVar(&path, "config", "", "path config path")
	flag.Parse()

	return path
}
