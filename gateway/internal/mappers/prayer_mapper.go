package mappers

import (
	"strings"

	prayersv1 "github.com/Jeudry/adventist-stack/gen/prayers/v1"
	"github.com/Jeudry/adventist-stack/gateway/internal/models/base"
	"github.com/Jeudry/adventist-stack/gateway/internal/models/prayer"
)

func ToPrayerVM(p *prayersv1.Prayer) prayer.PrayerVM {
	if p == nil {
		return prayer.PrayerVM{}
	}

	return prayer.PrayerVM{
		BaseVM:      base.ToBaseVM(p.GetId(), p.GetCreatedAt(), p.GetUpdatedAt()),
		Title:       p.GetTitle(),
		Description: p.GetDescription(),
		AuthorName:  p.GetAuthorName(),
		IsAnonymous: p.GetIsAnonymous(),
		Status:      p.GetStatus().String(),
	}
}

func ToCreatePrayerProto(req prayer.CreatePrayerRequest) *prayersv1.CreatePrayerRequest {
	return &prayersv1.CreatePrayerRequest{
		Title:       req.Title,
		Description: req.Description,
		AuthorName:  &req.AuthorName,
		IsAnonymous: req.IsAnonymous,
		Status:      PrayerStatusToProto(req.Status),
	}
}

func ToUpdatePrayerProto(id string, req prayer.UpdatePrayerRequest) *prayersv1.UpdatePrayerRequest {
	return &prayersv1.UpdatePrayerRequest{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
		AuthorName:  &req.AuthorName,
		IsAnonymous: req.IsAnonymous,
	}
}

func PrayerStatusToProto(s string) prayersv1.PrayerStatus {
	switch strings.ToLower(s) {
	case "pending":
		return prayersv1.PrayerStatus_PRAYER_STATUS_PENDING
	case "answered":
		return prayersv1.PrayerStatus_PRAYER_STATUS_ANSWERED
	case "archived":
		return prayersv1.PrayerStatus_PRAYER_STATUS_ARCHIVED
	default:
		return prayersv1.PrayerStatus_PRAYER_STATUS_UNSPECIFIED
	}
}
