package domain

import (
	"errors"
	"time"

	"github.com/Jeudry/adventist-stack/pkg/entity"
	"github.com/google/uuid"
)

type Enrollment struct {
	entity.Base
	SabbathSchoolId uuid.UUID
	MemberId        uuid.UUID
	EnrolledAt      time.Time
	UnenrolledAt    *time.Time
}

var (
	ErrEnrollmentNotFound      = errors.New("Enrollment not found")
	ErrInvalidEnrollment       = errors.New("Invalid enrollment")
	ErrInvalidUnenrolledDate   = errors.New("Invalid unenrolled date")
	ErrEnrollmentAlreadyClosed = errors.New("Enrollment already closed")
)

func (e *Enrollment) IsActive() bool {
	return e.UnenrolledAt == nil
}

func (e *Enrollment) Unenroll(at time.Time, actorId uuid.UUID) error {
	if e.IsActive() {
		return ErrEnrollmentAlreadyClosed
	}
	isInFuture := at.After(time.Now())
	isBeforeEnrollmentDate := at.Before(e.EnrolledAt)
	if isInFuture || isBeforeEnrollmentDate {
		return ErrInvalidUnenrolledDate
	}
	e.UpdatedAt = time.Now().UTC()
	e.UpdatedBy = &actorId
	e.UnenrolledAt = &at

	return nil
}

func (e *Enrollment) Normalize() {
	if e.EnrolledAt.IsZero() {
		e.EnrolledAt = time.Now().UTC().Truncate(24 * time.Hour)
	}
}

func (e Enrollment) Validate() error {
	return errors.Join(
		validateEnrollmentIdsRelation(e.SabbathSchoolId, e.MemberId),
		validateEnrolledDates(e.EnrolledAt, e.UnenrolledAt),
	)
}

func validateEnrolledDates(enrolledAt time.Time, unenrolledAt *time.Time) error {
	if enrolledAt.IsZero() {
		return ErrEnrollmentNotFound
	}
	if unenrolledAt != nil && unenrolledAt.Before(enrolledAt) {
		return ErrInvalidUnenrolledDate
	}
	return nil
}

func validateEnrollmentIdsRelation(sabbathSchoolId uuid.UUID, memberId uuid.UUID) error {
	if sabbathSchoolId == uuid.Nil || memberId == uuid.Nil {
		return ErrInvalidEnrollment
	}
	return nil
}
