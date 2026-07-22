// Package protoconv converts between protobuf wire types and Go domain types.
package protoconv

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// TimeFromProto converts an optional protobuf timestamp to *time.Time.
func TimeFromProto(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

// TimeToProto converts an optional *time.Time to a protobuf timestamp.
func TimeToProto(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
