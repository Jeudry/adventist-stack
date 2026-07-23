package clients

import (
	"fmt"

	"google.golang.org/grpc"

	membersv1 "github.com/Jeudry/adventist-stack/gen/members/v1"
)

func newMembersClient(addr string) (membersv1.MemberServiceClient, *grpc.ClientConn, error) {
	conn, err := dial(addr)
	if err != nil {
		return nil, nil, fmt.Errorf("clients: members dial: %w", err)
	}
	return membersv1.NewMemberServiceClient(conn), conn, nil
}
