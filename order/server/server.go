package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	account "github.com/wignn/micro-3/account/client"
	catalog "github.com/wignn/micro-3/catalog/client"
	"github.com/wignn/micro-3/order/genproto"
	"github.com/wignn/micro-3/order/model"
	"github.com/wignn/micro-3/order/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service       service.OrderService
	accountClient *account.AccountClient
	catalogClient *catalog.CatalogClient
	genproto.UnimplementedOrderServiceServer
}


func ListenGRPC(s service.OrderService, accountURL, catalogURL string, port int) error {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		accountClient.Close()
		return err
	}
	
	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		catalogClient.Close()	
		return err
	}

	serv := grpc.NewServer()
	genproto.RegisterOrderServiceServer(serv, &grpcServer{
		service:       s,
		accountClient: accountClient,
		catalogClient: catalogClient,
	})

	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostOrder(c context.Context, r *genproto.PostOrderRequest) (*genproto.PostOrderResponse, error) {
		// Check if account exists
	_, err := s.accountClient.GetAccount(c, r.AccountId)
	if err != nil {
		log.Println("Error getting account: ", err)
		return nil, errors.New("account not found")
	}

	// Get ordered products
	productIDs := []string{}
	for _, p := range r.Products {
		productIDs = append(productIDs, p.ProductId)
	}
	orderedProducts, err := s.catalogClient.GetProducts(c, 0, 0, productIDs, "")
	if err != nil {
		log.Println("Error getting products: ", err)
		return nil, errors.New("products not found")
	}

	// Construct products
	products := []model.OrderedProduct{}
	for _, p := range orderedProducts {
		product := model.OrderedProduct{
			ID:          p.Id,
			Quantity:    0,
			Price:       p.Price,
			Name:        p.Name,
			Description: p.Description,
		}
		for _, rp := range r.Products {
			if rp.ProductId == p.Id {
				product.Quantity = rp.Quantity
				break
			}
		}

		if product.Quantity != 0 {
			products = append(products, product)
		}
	}

	// Call service implementation
	order, err := s.service.PostOrder(c, r.AccountId, products)
	if err != nil {
		log.Println("Error posting order: ", err)
		return nil, errors.New("could not post order")
	}

	// Make response order
	orderProto := &genproto.Order{
		Id:         order.ID,
		AccountId:  order.AccountID,
		TotalPrice: order.TotalPrice,
		Products:   []*genproto.Order_OrderProduct{},
	}
	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()
	for _, p := range order.Products {
		orderProto.Products = append(orderProto.Products, &genproto.Order_OrderProduct{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		})
	}
	return &genproto.PostOrderResponse{
		Order: orderProto,
	}, nil
}


func (s *grpcServer) GetOrdersForAccount(
	ctx context.Context,
	r *genproto.GetOrdersForAccountRequest,
) (*genproto.GetOrdersForAccountResponse, error) {
	// Get orders for account
	accountOrders, err := s.service.GetOrdersForAccount(ctx, r.AccountId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Get all ordered products
	productIDMap := map[string]bool{}
	for _, o := range accountOrders {
		for _, p := range o.Products {
			productIDMap[p.ID] = true
		}
	}
	productIDs := []string{}
	for id := range productIDMap {
		productIDs = append(productIDs, id)
	}
	products, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("Error getting account products: ", err)
		return nil, err
	}

	// Construct orders
	orders := []*genproto.Order{}
	for _, o := range accountOrders {
		// Encode order
		op := &genproto.Order{
			AccountId:  o.AccountID,
			Id:         o.ID,
			TotalPrice: o.TotalPrice,
			Products:   []*genproto.Order_OrderProduct{},
		}
		op.CreatedAt, _ = o.CreatedAt.MarshalBinary()

		// Decorate orders with products
		for _, product := range o.Products {
			// Populate product fields
			for _, p := range products {
				if p.Id == product.ID {
					product.Name = p.Name
					product.Description = p.Description
					product.Price = p.Price
					break
				}
			}

			op.Products = append(op.Products, &genproto.Order_OrderProduct{
				Id:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Quantity:    product.Quantity,
			})
		}

		orders = append(orders, op)
	}
	return &genproto.GetOrdersForAccountResponse{Orders: orders}, nil
}


func (s *grpcServer) DeleteOrder(ctx context.Context, r *genproto.DeleteOrderRequest) (*genproto.DeleteOrderResponse, error) {
	err := s.service.DeleteOrder(ctx, r.Id)
	
	if err != nil {
		log.Println("Error deleting order: ", err)
		return nil, errors.New("could not delete order")
	}

	return &genproto.DeleteOrderResponse{
		DeletedId: r.Id,
		Message:   "Order deleted successfully",
		Success:   true,
	}, nil
}