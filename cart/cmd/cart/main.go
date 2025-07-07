package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"github.com/wignn/micro-3/cart/repository"
	"github.com/wignn/micro-3/cart/server"
	"github.com/wignn/micro-3/cart/service"
	
)

type Config struct {
	DSN  string `envconfig:"DATABASE_URL"`
	PORT int    `envconfig:"PORT" default:"50051"`
	CatalogURL  string `envconfig:"CATALOG_SERVICE_URL"`
	AccountURL  string `envconfig:"ACCOUNT_SERVICE_URL"`
}

func main() {
	var cfg Config
	
	fmt.Println("Starting Cart Service...")
	if err := envconfig.Process("", &cfg); err != nil {
		panic(err)
	}
	
	var r repository.CartRepository

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
	s := service.NewCartService(r)
	log.Fatal(server.ListenGRPC(s, cfg.CatalogURL, cfg.AccountURL, cfg.PORT))
}
