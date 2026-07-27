package clients

import (
	"fmt"

	prayersv1 "github.com/Jeudry/adventist-stack/gen/prayers/v1"
	"google.golang.org/grpc"
)

func newPrayersClient(addr string) (prayersv1.PrayerServiceClient, *grpc.ClientConn, error) {
	conn, err := dial(addr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial: %w", err)
	}
	return prayersv1.NewPrayerServiceClient(conn), conn, nil
}
