package models

import (
	"time"

	"github.com/google/uuid"
)

type Petition struct {
	ID          uuid.UUID
	CityID      uuid.UUID
	CreatorID   uuid.UUID
	Title       string
	Description string
	Status      string
	Signatures  int
	Goal        int
	Reply       string
	EndDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PetitionSignature struct {
	ID         uuid.UUID
	PetitionID uuid.UUID
	UserID     uuid.UUID
	CreatedAt  time.Time
}
