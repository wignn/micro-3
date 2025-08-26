package main

import (
    "log"
    "net/http"

    "github.com/99designs/gqlgen/graphql/handler"
    "github.com/99designs/gqlgen/graphql/playground"
    "github.com/gorilla/handlers"
    "github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
    AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
    CatalogURL string `envconfig:"CATALOG_SERVICE_URL"`
    OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
    ReviewURL  string `envconfig:"REVIEW_SERVICE_URL"`
    AuthURL    string `envconfig:"AUTH_SERVICE_URL"`
    CartURL    string `envconfig:"CART_SERVICE_URL"`
}

func main() {
    var cfg AppConfig

    err := envconfig.Process("", &cfg)
    if err != nil {
        log.Fatalf("failed to process env config: %v", err)
    }

    s, err := NewGraphQLServer(
        cfg.AccountURL,
        cfg.CatalogURL,
        cfg.OrderURL,
        cfg.ReviewURL,
        cfg.AuthURL,
        cfg.CartURL,
    )
    
    if err != nil {
        log.Fatalf("failed to create GraphQL server: %v", err)
    }

    schema, err := s.ToExecutableSchema()
    if err != nil {
        log.Fatalf("failed to create schema: %v", err)
    }


    mux := http.NewServeMux()
    mux.Handle("/graphql", handler.NewDefaultServer(schema))
    mux.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))
    corsHandler := handlers.CORS(
        handlers.AllowedOrigins([]string{
            "http://localhost:3000",   
            "http://54.242.132.124:3000",  
        }),
        handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
        handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
        handlers.AllowCredentials(),
    )(mux)

    log.Println("Server running at http://localhost:8000")
    log.Fatal(http.ListenAndServe(":8000", corsHandler))
}
