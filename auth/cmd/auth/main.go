package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"github.com/wignn/micro-3/auth/repository"
	"github.com/wignn/micro-3/auth/server"
	"github.com/wignn/micro-3/auth/service"
	"github.com/wignn/micro-3/auth/utils"
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

	fmt.Printf("Starting Review Service on port %d...\n", cfg.PORT)
	fmt.Printf("Using database DSN: %s\n", cfg.DSN)
	var r repository.AuthRepository

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = repository.NewAuthPostgresRepository(cfg.DSN)
		if err != nil {
			log.Println("Failed to connect to database, retrying...")
			return err
		}
		return
	})
	utils.InitJWTConfig()
	defer r.Close()
	log.Println("listening on port", cfg.PORT)
	s := service.NewAuthService(r)
	log.Fatal(server.ListenGRPC(s, cfg.PORT))
}
