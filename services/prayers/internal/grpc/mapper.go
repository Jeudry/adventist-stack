package grpc

import (
	"errors"

	prayersv1 "github.com/Jeudry/adventist-stack/gen/prayers/v1"
	"github.com/Jeudry/adventist-stack/services/prayers/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func prayerFromCreateRequest(req *prayersv1.CreatePrayerRequest) (domain.Prayer, error) {
	p := domain.Prayer{
		Title:       req.Title,
		Description: req.Description,
		AuthorName:  *req.AuthorName,
		IsAnonymous: req.IsAnonymous,
		Status:      domain.Status(req.Status),
	}
	return p, nil
}

func prayerFromUpdateRequest(req *prayersv1.UpdatePrayerRequest) (domain.Prayer, error) {
	p := domain.Prayer{
		Title:       req.Title,
		Description: req.Description,
		AuthorName:  *req.AuthorName,
		IsAnonymous: req.IsAnonymous,
		Status:      domain.Status(req.Status),
	}
	return p, nil
}

func prayerToProto(p domain.Prayer) *prayersv1.Prayer {
	prayer := &prayersv1.Prayer{
		Title:       p.Title,
		Description: p.Description,
		AuthorName:  &p.AuthorName,
		IsAnonymous: p.IsAnonymous,
		Status:      prayersv1.PrayerStatus(p.Status),
		Id:          p.ID.String(),
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}
	return prayer
}

func statusFromProto(s prayersv1.PrayerStatus) domain.Status {
	switch s {
	case prayersv1.PrayerStatus_PRAYER_STATUS_PENDING:
		return domain.STATUS_PENDING
	case prayersv1.PrayerStatus_PRAYER_STATUS_ANSWERED:
		return domain.STATUS_ANSWERED
	case prayersv1.PrayerStatus_PRAYER_STATUS_ARCHIVED:
		return domain.STATUS_ARCHIVED
	default:
		return domain.STATUS_UNSPECIFIED
	}
}

func statusToProto(s domain.Status) prayersv1.PrayerStatus {
	switch s {
	case domain.STATUS_PENDING:
		return prayersv1.PrayerStatus_PRAYER_STATUS_PENDING
	case domain.STATUS_ANSWERED:
		return prayersv1.PrayerStatus_PRAYER_STATUS_ANSWERED
	case domain.STATUS_ARCHIVED:
		return prayersv1.PrayerStatus_PRAYER_STATUS_ARCHIVED
	default:
		return prayersv1.PrayerStatus_PRAYER_STATUS_UNSPECIFIED
	}
}

func toStatus(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidPrayer):
		return domain.ErrInvalidPrayer
	case errors.Is(err, domain.ErrInvalidStatus):
		return domain.ErrInvalidStatus
	case errors.Is(err, domain.ErrPrayerNotFound):
		return domain.ErrPrayerNotFound
	default:
		return err
	}
}
