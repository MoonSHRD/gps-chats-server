package config

import (
	"fmt"
	"os"
	"regexp"

	"github.com/MoonSHRD/logger"
	"github.com/pelletier/go-toml"
)

type Config struct {
	MongoDB       MongoDB `toml:"mongoDB"`
	HTTP          HTTP    `toml:"http"`
	AuthServerURL string  `toml:"authServerURL"`
	JWT           JWT     `toml:"jwt"`
}

type MongoDB struct {
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

type JWT struct {
	SigningKey string `toml:"signingKey"`
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
	err = validateConfig(&cfg)
	return cfg, err
}

func validateConfig(config *Config) error {
	matches, err := regexp.MatchString("http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\\(\\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+", config.AuthServerURL)
	if err != nil {
		return err
	}
	if !matches {
		return fmt.Errorf("invalid auth server url")
	}
	return nil
}
