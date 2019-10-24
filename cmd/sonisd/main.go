package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func main() {
	logger.Info("Starting microservice...")

	logger.Info("Microservice successfully started!")

	// CTRL+C handler.
	signalHandler := make(chan os.Signal, 1)
	shutdownDone := make(chan bool, 1)
	signal.Notify(signalHandler, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalHandler
		logger.Info("CTRL+C or SIGTERM received, shutting down openkeepd...")
		// TODO make graceful shutdown of microservice
		shutdownDone <- true
	}()

	<-shutdownDone
	logger.Info("Microservice successfully shutted down")
	os.Exit(0)
}
