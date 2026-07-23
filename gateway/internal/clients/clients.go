package clients

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "github.com/Jeudry/adventist-stack/gen/auth/v1"
	membersv1 "github.com/Jeudry/adventist-stack/gen/members/v1"
	notificationsv1 "github.com/Jeudry/adventist-stack/gen/notifications/v1"
)

type Config struct {
	AuthAddr          string
	NotificationsAddr string
	MembersAddr       string
}

type Clients struct {
	Auth          authv1.AuthServiceClient
	Notifications notificationsv1.NotificationServiceClient
	Members       membersv1.MemberServiceClient

	conns []*grpc.ClientConn
}

func New(cfg Config) (*Clients, error) {
	authClient, authConn, err := newAuthClient(cfg.AuthAddr)
	if err != nil {
		return nil, err
	}

	notifClient, notifConn, err := newNotificationsClient(cfg.NotificationsAddr)
	if err != nil {
		authConn.Close()
		return nil, err
	}

	memClient, memConn, err := newMembersClient(cfg.MembersAddr)
	if err != nil {
		authConn.Close()
		notifConn.Close()
		return nil, err
	}

	return &Clients{
		Auth:          authClient,
		Notifications: notifClient,
		Members:       memClient,
		conns:         []*grpc.ClientConn{authConn, notifConn, memConn},
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
