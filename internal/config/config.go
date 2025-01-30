package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env"`
	StoragePath string `yaml:"storage_path"`

	HTTPServer HTTPServer `yaml:"http_server"`
	Database   Database   `yaml:"database"`

	// TODO: Add the remaining structures
}
type HTTPServer struct {
	Address     string `yaml:"address"`
	Timeout     string `yaml:"timeout"`
	IdleTimeout string `yaml:"idle_timeout"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
}
type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}
type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
type Nats struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func MustLoad() *Config {
	// Указываем полный путь к конфигурационному файлу
	configPath := filepath.Join("../config/config.yaml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
