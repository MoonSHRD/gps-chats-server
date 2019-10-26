package main

import (
	"os"
	"os/signal"
	"syscall"

	internalCtx "github.com/MoonSHRD/sonis/internal/context"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()
var context *internalCtx.Context

func main() {
	var err error
	logger.Info("Starting microservice...")
	context, err = internalCtx.New()
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
