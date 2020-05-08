package config

import (
	"os"

	"github.com/MoonSHRD/logger"
	"github.com/pelletier/go-toml"
)

type Config struct {
	PostgreSQL PostgreSQL `toml:"postgresql"`
	HTTP       HTTP       `toml:"http"`
}

type PostgreSQL struct {
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	User         string `toml:"user"`
	Password     string `toml:"password"`
	DatabaseName string `toml:"databaseName"`
}

type HTTP struct {
	Address string `toml:"address"`
	Port    int    `toml:"port"`
}

func NewConfig(path string) (Config, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		logger.Fatal(err)
	}
	decoder := toml.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		logger.Fatal(err)
	}

	return cfg, nil
}
