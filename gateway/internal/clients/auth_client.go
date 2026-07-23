package clients

import (
	"fmt"

	"google.golang.org/grpc"

	authv1 "github.com/Jeudry/adventist-stack/gen/auth/v1"
)

func newAuthClient(addr string) (authv1.AuthServiceClient, *grpc.ClientConn, error) {
	conn, err := dial(addr)
	if err != nil {
		return nil, nil, fmt.Errorf("clients: auth dial: %w", err)
	}
	return authv1.NewAuthServiceClient(conn), conn, nil
}
