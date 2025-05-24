package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AccountUrl string `envconfig:"ACCOUNT_URL" envconfigDefault:"http://localhost:8081"`
	CatalogUrl string `envconfig:"CATALOG_URL" envconfigDefault:"http://localhost:8082"`
	OrderUrl   string `envconfig:"ORDER_URL" envconfigDefault:"http://localhost:8083"`
}

func main() {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("failed to process env config: %v", err)
	}

	s, err := NewGraphQLServer(cfg.AccountUrl, cfg.CatalogUrl, cfg.OrderUrl)
	
	if err != nil {
		log.Fatalf("failed to create GraphQL server: %v", err)
	}

	schema, err := s.ToExecutableSchema()
	
	if err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}

	http.Handle("/graphql", handler.New(schema))
	http.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}