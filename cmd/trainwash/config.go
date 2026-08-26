package main

import (
	"flag"
	"os"
	"path/filepath"
)

type Config struct {
	ListenAddr string
	DataDir    string
}

func LoadConfig() Config {
	cfg := Config{
		ListenAddr: envOr("TRAINWASH_ADDR", ":8080"),
		DataDir:    envOr("TRAINWASH_DATA", defaultDataDir()),
	}
	flag.StringVar(&cfg.ListenAddr, "addr", cfg.ListenAddr, "http listen address")
	flag.StringVar(&cfg.DataDir, "data", cfg.DataDir, "state data directory")
	flag.Parse()
	return cfg
}

func defaultDataDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "data"
	}
	return filepath.Join(dir, "data")
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
