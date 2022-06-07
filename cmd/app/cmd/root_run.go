package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/fancar/cobra_template/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func run(cnd *cobra.Command, args []string) error {
	// setLogLevel()
	tasks := []func(context.Context, *sync.WaitGroup) error{
		setLogLevel,
		printStartMessage,
		// startSomeRoutine,
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
		log.WithField("signal", s).Info("signal received, hard terminated")
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
	}).Info("starting %s ...", appName)
	return nil
}
