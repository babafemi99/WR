package main

import (
	"fmt"
	"github.com/babafemi99/WR/internal/config"
	"github.com/babafemi99/WR/internal/logger"
	"github.com/babafemi99/WR/pkg/deps"
	api "github.com/babafemi99/WR/pkg/transport/http/rest"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	allowConnectionsAfterShutdown = 5 * time.Second
)

func main() {
	cfg := config.New()
	logger.Init()

	nDeps := deps.New(cfg)

	a := api.API{
		Config: cfg,
		Deps:   nDeps,
	}

	go func() {
		logger.Log.Info(fmt.Sprintf("Server running on port %v ...", cfg.Port))
		log.Fatal(a.Serve())
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan

	logger.Log.Info("Request to shutdown server. Doing nothing for ", allowConnectionsAfterShutdown)
	waitTimer := time.NewTimer(allowConnectionsAfterShutdown)
	<-waitTimer.C

	logger.Log.Info("Shutting down server...")
	log.Fatal(a.Shutdown())
}
