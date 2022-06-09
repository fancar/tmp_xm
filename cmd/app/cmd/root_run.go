package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/fancar/tmp_xm/internal/api"
	"github.com/fancar/tmp_xm/internal/config"
	"github.com/fancar/tmp_xm/internal/storage"
)

func run(cnd *cobra.Command, args []string) error {
	tasks := []func(context.Context, *sync.WaitGroup) error{
		setLogLevel,
		printStartMessage,
		setupStorage,
		setupAPI,
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	for _, t := range tasks {
		if err := t(ctx, &wg); err != nil {
			log.Fatal(err)
		}
	}

	exitChan := make(chan struct{})
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	go func() {
		cancel()
		log.Info("Stopping gracefully ...")
		wg.Wait()
		log.Info("Bye!")
		exitChan <- struct{}{}
	}()
	cancel()
	select {
	case <-exitChan:
	case s := <-sigChan:
		log.WithField("signal", s).Info("signal received, terminated")
	}

	return nil
}

func setLogLevel(ctx context.Context, wg *sync.WaitGroup) error {
	log.SetLevel(log.Level(uint8(config.C.General.LogLevel)))
	return nil
}

func printStartMessage(ctx context.Context, wg *sync.WaitGroup) error {
	log.WithFields(log.Fields{
		"version": version,
		// "docs":    "https://www. ... .su/",
	}).Infof("starting %s ...", appName)
	return nil
}

func setupStorage(ctx context.Context, wg *sync.WaitGroup) error {
	if err := storage.Setup(config.C); err != nil {
		return fmt.Errorf("can't setup storage: %v", err)
	}

	return nil
}

func setupAPI(ctx context.Context, wg *sync.WaitGroup) error {
	if err := api.Setup(ctx, config.C); err != nil {
		return fmt.Errorf("can't setup api: %v", err)
	}
	return nil
}
