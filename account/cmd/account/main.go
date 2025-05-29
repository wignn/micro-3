package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"github.com/wignn/micro-3/account/repository"
	"github.com/wignn/micro-3/account/server"
	"github.com/wignn/micro-3/account/service"
)

type config struct {
	Port int    `env:"PORT" envDefault:"50051"`
	Host string `env:"HOST" envDefault:"localhost"`
	DSN  string `env:"DSN"`
}

func main() {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal("Failed to process environment variables:", err)
	}

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

	log.Println("listening on port", cfg.Port)
	s := service.NewAccountService(r)
	log.Fatal(server.ListenGRPC(s, cfg.Port))	
}
