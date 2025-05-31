package main

import (
	"log"
	"time"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"fmt"
	"github.com/wignn/micro-3/account/repository"
	"github.com/wignn/micro-3/account/server"
	"github.com/wignn/micro-3/account/service"
)

type Config struct {
	DSN  string `envconfig:"DATABASE_URL"`
	PORT int    `envconfig:"PORT" default:"50051"`
}

func main() {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal("Failed to process environment variables:", err)
	}
	fmt.Println("Starting Account Service...", cfg.PORT)

	var r repository.AccountRepository

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = repository.NewPostgresRepository(cfg.DSN)
		if err != nil {
			log.Println("Failed to connect to database, retrying...")
			return err
		}
		return nil
	})
	defer r.Close()

	log.Println("listening on port", cfg.PORT)
	s := service.NewAccountService(r)
	log.Fatal(server.ListenGRPC(s, cfg.PORT))
}
