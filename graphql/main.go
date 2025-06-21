package main

import (
	"log"
	"net/http"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
	ReviewURL  string `envconfig:"REVIEW_SERVICE_URL"`
}

func main() {
	var cfg AppConfig
	
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("failed to process env config: %v", err)
	}

	s, err := NewGraphQLServer(cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL, cfg.ReviewURL)
	if err != nil {
		log.Fatalf("failed to create GraphQL server: %v", err)
	}

	schema, err := s.ToExecutableSchema()
	
	if err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}

	http.Handle("/graphql", handler.NewDefaultServer(schema))
	http.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}