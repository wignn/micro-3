package main

import (
	"github.com/99designs/gqlgen/graphql"
	account "github.com/wignn/micro-3/account/client"
	auth "github.com/wignn/micro-3/auth/client"
	catalog "github.com/wignn/micro-3/catalog/client"
	order "github.com/wignn/micro-3/order/client"
	review "github.com/wignn/micro-3/review/client"
	cart "github.com/wignn/micro-3/cart/client"
)

type GraphQLServer struct {
	accountClient *account.AccountClient
	catalogClient *catalog.CatalogClient
	orderClient   *order.OrderClient
	authClient    *auth.AuthClient
	reviewClient  *review.ReviewClient
	cartClient    *cart.CartClient
}

func NewGraphQLServer(accountUrl, catalogUrl, orderUrl, reviewUrl, authUrl, cartUrl string) (*GraphQLServer, error) {
	accoutClient, err := account.NewClient(accountUrl)
	if err != nil {
		return nil, err
	}

	catalogClient, err := catalog.NewClient(catalogUrl)
	if err != nil {
		accoutClient.Close()
		return nil, err
	}

	orderClient, err := order.NewClient(orderUrl)
	if err != nil {
		catalogClient.Close()
		return nil, err
	}

	reviewClient, err := review.NewClient(reviewUrl)
	if err != nil {
		reviewClient.Close()
		return nil, err
	}

	authClient, err := auth.NewClient(authUrl)
	if err != nil {
		reviewClient.Close()
		return nil, err
	}

	cartClient, err := cart.NewClient(cartUrl)
	if err != nil {
		authClient.Close()
		reviewClient.Close()
		return nil, err
	}
	
	return &GraphQLServer{
		accoutClient,
		catalogClient,
		orderClient,
		authClient,
		reviewClient,
		cartClient,
	}, nil
}

func (s *GraphQLServer) Mutation() MutationResolver {
	return &mutationResolver{
		server: s,
	}
}

func (s *GraphQLServer) Query() QueryResolver {
	return &queryResolver{
		server: s,
	}
}

func (s *GraphQLServer) Account() AccountResolver {
	return &accountResolver{
		server: s,
	}
}

func (s *GraphQLServer) ToExecutableSchema() (graphql.ExecutableSchema, error) {
	return NewExecutableSchema(Config{
		Resolvers: s,
	}), nil
}
