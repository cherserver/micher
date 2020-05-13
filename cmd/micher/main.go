package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/cherserver/micher/service/application"
	"github.com/cherserver/micher/service/config"
	"github.com/cherserver/micher/service/environment"
)

func processCmdAndGetConfigPath() string {
	rootCmd := &cobra.Command{
		Use:     "micher",
		Short:   "Xiaomi hub API gateway",
		Long:    `MiCher is HTTP API gateway for Xiaomi hub and all of its connected devices`,
		Version: "0.1",
	}

	err := rootCmd.Execute()

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to run micher:", err)
		os.Exit(1)
	}

	return ""
}

func waitForExit(shutdownCtx context.Context) {
	signalCh := make(chan os.Signal, 1)
	done := make(chan struct{}, 1)

	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case gotSignal := <-signalCh:
			log.Printf("Got signal '%v', exit", gotSignal)
		case <-shutdownCtx.Done():
			log.Printf("Context shut down, exit")
		}

		done <- struct{}{}
	}()

	log.Println("awaiting signal")
	<-done
	log.Println("exiting")
}

func main() {
	configPath := processCmdAndGetConfigPath()

	var err error
	cfg := config.New(configPath)

	var env = environment.New(cfg)
	if err = env.Init(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to initialize environment:", err)
		os.Exit(1)
	}

	app := application.New(env)
	if err = app.Init(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to initialize application:", err)
		os.Exit(1)
	}

	waitForExit(env.ShutdownCtx())

}
