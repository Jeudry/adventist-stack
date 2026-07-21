package clients

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "github.com/Jeudry/adventist-stack/gen/auth/v1"
	notificationsv1 "github.com/Jeudry/adventist-stack/gen/notifications/v1"
)

type Clients struct {
	Auth          authv1.AuthServiceClient
	Notifications notificationsv1.NotificationServiceClient

	conns []*grpc.ClientConn
}

func New(authAddr, notificationsAddr string) (*Clients, error) {
	authConn, err := dial(authAddr)
	if err != nil {
		return nil, fmt.Errorf("clients: auth: %w", err)
	}

	notifConn, err := dial(notificationsAddr)
	if err != nil {
		authConn.Close()
		return nil, fmt.Errorf("clients: notifications: %w", err)
	}

	return &Clients{
		Auth:          authv1.NewAuthServiceClient(authConn),
		Notifications: notificationsv1.NewNotificationServiceClient(notifConn),
		conns:         []*grpc.ClientConn{authConn, notifConn},
	}, nil
}

func (c *Clients) Close() {
	for _, conn := range c.conns {
		_ = conn.Close()
	}
}

func dial(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
