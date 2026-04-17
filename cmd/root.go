package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/isyuricunha/gitflect/internal/config"
	"github.com/isyuricunha/gitflect/internal/syncer"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitflect",
	Short: "Automatically mirrors git repositories between providers",
	RunE:  run,
}

func Execute() error {
	return rootCmd.Execute()
}

func run(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Starting gitflect", "source", cfg.SourceProvider, "dest", cfg.DestProvider)

	s := syncer.New(cfg, logger)

	logger.Info("Starting initial synchronization loop")
	s.Run(context.Background())

	// If no sync interval is configured, run only once and then exit.
	if cfg.SyncInterval == "" {
		logger.Info("No SYNC_INTERVAL defined, finishing execution.")
		return nil
	}

	c := cron.New()
	schedule := fmt.Sprintf("@every %s", cfg.SyncInterval)
	if _, err := c.AddFunc(schedule, func() {
		s.Run(context.Background())
	}); err != nil {
		return fmt.Errorf("failed to schedule cron %s: %w", schedule, err)
	}
	
	c.Start()
	logger.Info("Scheduler started", "interval", cfg.SyncInterval)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	logger.Info("Shutting down gitflect safely...")
	c.Stop()
	return nil
}
