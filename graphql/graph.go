package main

import (
	"github.com/99designs/gqlgen/graphql"
	account "github.com/wignn/micro-3/account/client"
	catalog "github.com/wignn/micro-3/catalog/client"
	order "github.com/wignn/micro-3/order/client"
)

type GraphQLServer struct {
	accountClient *account.AccountClient
	catalogClient *catalog.CatalogClient
	orderClient  *order.Client
}

func NewGraphQLServer(accountUrl, catalogUrl, orderUrl string) (*GraphQLServer, error) {

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

	return &GraphQLServer{
		accoutClient,
		catalogClient,
		orderClient,
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