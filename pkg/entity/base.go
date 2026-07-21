// Package entity holds building blocks shared by every domain entity.
package entity

import (
	"time"

	"github.com/google/uuid"
)

// Base carries the identity and audit fields every persisted entity shares.
// Embed it (Go composition) so its fields are promoted onto the entity:
//
//	type Member struct {
//		entity.Base
//		FirstName string
//	}
//
// then use member.ID, member.CreatedAt, member.IsDeleted() directly.
type Base struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	CreatedBy *uuid.UUID
	UpdatedBy *uuid.UUID
	DeletedBy *uuid.UUID
}

// IsDeleted reports whether the entity has been soft-deleted.
func (b Base) IsDeleted() bool { return b.DeletedAt != nil }
