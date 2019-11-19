package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	internalCtx "github.com/MoonSHRD/sonis/internal/context"
	"github.com/MoonSHRD/sonis/internal/utils"
	"gopkg.in/yaml.v2"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()
var context *internalCtx.Context

func main() {
	var err error
	logger.Info("Starting microservice...")

	var cfg utils.Config
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to config")
	flag.Parse()
	if configPath == "" {
		logger.Error("Path to config isn't specified!")
		os.Exit(1)
	}
	cfgData, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.Error("Failed to read config!")
		os.Exit(1)
	}
	err = yaml.Unmarshal(cfgData, &cfg)
	if err != nil {
		logger.Errorf("Failed to read config! (yaml error: %s)", err.Error())
		os.Exit(1)
	}

	context, err = internalCtx.New(cfg)
	if err != nil {
		os.Exit(1)
	}
	logger.Info("Microservice successfully started!")

	// CTRL+C handler.
	signalHandler := make(chan os.Signal, 1)
	shutdownDone := make(chan bool, 1)
	signal.Notify(signalHandler, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalHandler
		logger.Info("CTRL+C or SIGTERM received, shutting down sonisd...")
		context.Destroy()
		shutdownDone <- true
	}()

	<-shutdownDone
	logger.Info("Microservice successfully shutted down")
	os.Exit(0)
}
