package clients

import (
	"fmt"

	"google.golang.org/grpc"

	notificationsv1 "github.com/Jeudry/adventist-stack/gen/notifications/v1"
)

func newNotificationsClient(addr string) (notificationsv1.NotificationServiceClient, *grpc.ClientConn, error) {
	conn, err := dial(addr)
	if err != nil {
		return nil, nil, fmt.Errorf("clients: notifications dial: %w", err)
	}
	return notificationsv1.NewNotificationServiceClient(conn), conn, nil
}
