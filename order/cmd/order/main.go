package main

import (
	"log"
	"time"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"github.com/wignn/micro-3/order/repository"
	"github.com/wignn/micro-3/order/server"
	"github.com/wignn/micro-3/order/service"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
	AccountURL  string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL  string `envconfig:"CATALOG_SERVICE_URL"`
	PORT        int    `envconfig:"PORT" default:"50051"`
}

func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal("Failed to process environment variables:", err)
	}

	var r repository.OrderRepository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = repository.NewOrderPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()

	s := service.NewOrderService(r)

	log.Println("Order service listening on port", cfg.PORT)
	log.Fatal(server.ListenGRPC(s, cfg.AccountURL, cfg.CatalogURL, cfg.PORT))
}
