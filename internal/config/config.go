package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	SourceProvider string
	SourceToken    string
	SourceUser     string

	DestProvider string
	DestToken    string
	DestUser     string
	DestURL      string 

	Visibility   string
	Include      []string
	Exclude      []string

	SyncInterval string
}

func Load() (*Config, error) {
	cfg := &Config{
		SourceProvider: env("SOURCE_PROVIDER", "github"),
		SourceToken:    os.Getenv("SOURCE_TOKEN"),
		SourceUser:     os.Getenv("SOURCE_USER"),
		DestProvider:   env("DEST_PROVIDER", "gitlab"),
		DestToken:      os.Getenv("DEST_TOKEN"),
		DestUser:       os.Getenv("DEST_USER"),
		DestURL:        os.Getenv("DEST_URL"),
		Visibility:     env("REPO_VISIBILITY", "all"),
		SyncInterval:   os.Getenv("SYNC_INTERVAL"),
	}

	if v := os.Getenv("REPO_INCLUDE"); v != "" {
		cfg.Include = splitTrim(v)
	}
	if v := os.Getenv("REPO_EXCLUDE"); v != "" {
		cfg.Exclude = splitTrim(v)
	}

	required := map[string]string{
		"SOURCE_TOKEN": cfg.SourceToken,
		"SOURCE_USER":  cfg.SourceUser,
		"DEST_TOKEN":   cfg.DestToken,
		"DEST_USER":    cfg.DestUser,
	}
	for k, v := range required {
		if v == "" {
			return nil, fmt.Errorf("%s is required", k)
		}
	}

	return cfg, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func splitTrim(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
