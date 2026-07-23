package clients

import (
	"fmt"

	"google.golang.org/grpc"

	productsv1 "github.com/Jeudry/adventist-stack/gen/products/v1"
)

func newProductsClient(addr string) (productsv1.ProductServiceClient, *grpc.ClientConn, error) {
	conn, err := dial(addr)
	if err != nil {
		return nil, nil, fmt.Errorf("clients: products dial: %w", err)
	}
	return productsv1.NewProductServiceClient(conn), conn, nil
}
