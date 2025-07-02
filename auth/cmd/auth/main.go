package main

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DSN  string `envconfig:"DATABASE_URL"`
	PORT int    `envconfig:"PORT" default:"50051"`
}

func main() {
	var cfg Config
	fmt.Println("Starting Auth Service...")

	if err := envconfig.Process("", &cfg); err != nil {
		fmt.Println("Failed to process environment variables:", err)
	}

	fmt.Printf("Starting Auth Service on port %d...\n", cfg.PORT)

	log.Println("listening on port", cfg.PORT)
	s := 
}