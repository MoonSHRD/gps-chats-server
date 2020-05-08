package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/MoonSHRD/logger"
	"github.com/pelletier/go-toml"

	"github.com/MoonSHRD/sonis/app"
	"github.com/MoonSHRD/sonis/config"
	"github.com/MoonSHRD/sonis/router"
)

func main() {
	var configPath string
	var generateConfig bool
	var verboseLogging bool
	var syslogLogging bool
	flag.StringVar(&configPath, "config", "", "Path to server config")
	flag.BoolVar(&generateConfig, "genconfig", false, "Generate new config")
	flag.BoolVar(&verboseLogging, "verbose", true, "Verbose logging")
	flag.BoolVar(&syslogLogging, "syslog", false, "Log to system logging daemon")
	flag.Parse()
	defer logger.Init("sonis", verboseLogging, syslogLogging, ioutil.Discard).Close() // TODO Make ability to use file for log output
	if generateConfig {
		config := config.Config{}
		configStr, err := toml.Marshal(config)
		if err != nil {
			logger.Fatalf("Cannot generate config! %s", err.Error())
		}
		fmt.Print(string(configStr))
		os.Exit(0)
	}
	logger.Info("Starting Sonis...")
	if configPath == "" {
		logger.Fatal("Path to config isn't specified!")
	}

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	app, err := app.NewApp(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	router, err := router.NewRouter(app)
	if err != nil {
		logger.Fatalf("Failed to initialize router: %s", err.Error())
	}
	app.Run(router)
}
