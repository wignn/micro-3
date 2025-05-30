package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"github.com/wignn/micro-3/catalog/repository"
	"github.com/wignn/micro-3/catalog/server"
	"github.com/wignn/micro-3/catalog/service"
)

type Config struct {
	DSN string `envconfig:"DATABASE_URL"`
	PORT int    `envconfig:"PORT" default:"50051"`
}


func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal("Failed to process environment variables:", err)
	}
	
	var r repository.CatalogRepository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = repository.NewElasticRepository(cfg.DSN)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()

	log.Println("Listening on port", cfg.PORT)
	s := service.NewCatalogService(r)
	log.Fatal(server.ListenGRPC(s, cfg.PORT))
}