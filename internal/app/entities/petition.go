package entities

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/city-petitions-svc/internal/app/models"
	"github.com/chains-lab/city-petitions-svc/internal/constant/enum"
	"github.com/chains-lab/city-petitions-svc/internal/dbx"
	"github.com/chains-lab/city-petitions-svc/internal/errx"
	"github.com/chains-lab/city-petitions-svc/internal/pagination"
	"github.com/google/uuid"
)

type petitionsQ interface {
	New() dbx.PetitionsQ

	Insert(ctx context.Context, in dbx.Petition) error
	Get(ctx context.Context) (dbx.Petition, error)
	Select(ctx context.Context) ([]dbx.Petition, error)
	Update(ctx context.Context, in dbx.UpdatePetitionInput) error
	Delete(ctx context.Context) error

	FilterID(id uuid.UUID) dbx.PetitionsQ
	FilterCityID(cityID uuid.UUID) dbx.PetitionsQ
	FilterCreatorID(userID uuid.UUID) dbx.PetitionsQ
	FilterStatus(status string) dbx.PetitionsQ

	FilterStatusIn(statuses ...string) dbx.PetitionsQ

	FilterCreatedAt(t time.Time, after bool) dbx.PetitionsQ
	FilterEndDate(t time.Time, after bool) dbx.PetitionsQ

	TitleLike(s string) dbx.PetitionsQ

	OrderByCreated(ascending bool) dbx.PetitionsQ
	OrderBySignatures(ascending bool) dbx.PetitionsQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.PetitionsQ
}

type signaturesQ interface {
	New() dbx.PetitionSignaturesQ

	Insert(ctx context.Context, input dbx.PetitionSignature) error
	Get(ctx context.Context) (dbx.PetitionSignature, error)
	Select(ctx context.Context) ([]dbx.PetitionSignature, error)
	Delete(ctx context.Context) error

	FilterID(id uuid.UUID) dbx.PetitionSignaturesQ
	FilterPetitionID(petitionID uuid.UUID) dbx.PetitionSignaturesQ
	FilterUserID(userID uuid.UUID) dbx.PetitionSignaturesQ

	OrderByCreated(ascending bool) dbx.PetitionSignaturesQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.PetitionSignaturesQ
}

type Petition struct {
	q    petitionsQ
	sigQ signaturesQ
}

func NewPetition(pg *sql.DB) Petition {
	return Petition{
		q:    dbx.NewPetitionsQ(pg),
		sigQ: dbx.NewPetitionSignaturesQ(pg),
	}
}

type CreatePetitionInput struct {
	Title       string
	Description string
}

func (p Petition) CreatePetition(ctx context.Context, cityID, creatorID uuid.UUID, input CreatePetitionInput) (models.Petition, error) {
	petitionID := uuid.New()
	now := time.Now().UTC()

	petition := dbx.Petition{
		ID:          petitionID,
		CityID:      cityID,
		CreatorID:   creatorID,
		Title:       input.Title,
		Description: input.Description,
		Status:      enum.PetitionPublished,
		Signatures:  0,
		Goal:        10000, // Default goal is 10,000 signatures
		Reply:       "",
		EndDate:     now.AddDate(0, 0, 30), // Default end date is 30 days from now
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := p.q.New().Insert(ctx, petition); err != nil {
		switch {
		default:
			return models.Petition{}, errx.RaiseInternal(ctx, err)
		}
	}

	return petitionModel(petition), nil
}

func (p Petition) GetPetition(ctx context.Context, petitionID uuid.UUID) (models.Petition, error) {
	petition, err := p.q.New().FilterID(petitionID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Petition{}, errx.RaisePetitionNotFoundByID(ctx, err, petitionID)
		default:
			return models.Petition{}, errx.RaiseInternal(ctx, err)
		}
	}

	return petitionModel(petition), nil
}

func (p Petition) ApprovePetition(ctx context.Context, petitionID uuid.UUID, reply string) (models.Petition, error) {
	petition, err := p.q.New().FilterID(petitionID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Petition{}, errx.RaisePetitionNotFoundByID(ctx, err, petitionID)
		default:
			return models.Petition{}, errx.RaiseInternal(ctx, err)
		}
	}

	status := enum.PetitionApproved

	updateInput := dbx.UpdatePetitionInput{
		Status: &status,
		Reply:  &reply,
	}

	if err := p.q.New().FilterID(petitionID).Update(ctx, updateInput); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Petition{}, errx.RaisePetitionNotFoundByID(ctx, err, petitionID)
		default:
			return models.Petition{}, errx.RaiseInternal(ctx, err) // Other errors
		}
	}

	//TODO add kafka event for petition approval

	return models.Petition{
		ID:          petition.ID,
		CityID:      petition.CityID,
		CreatorID:   petition.CreatorID,
		Title:       petition.Title,
		Description: petition.Description,
		Status:      status,
		Signatures:  petition.Signatures,
		Goal:        petition.Goal,
		Reply:       reply,
		EndDate:     petition.EndDate,
		CreatedAt:   petition.CreatedAt,
		UpdatedAt:   petition.UpdatedAt,
	}, nil
}

func (p Petition) RejectPetition(ctx context.Context, petitionID uuid.UUID, reply string) (models.Petition, error) {
	petition, err := p.q.New().FilterID(petitionID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Petition{}, errx.RaisePetitionNotFoundByID(ctx, err, petitionID)
		default:
			return models.Petition{}, errx.RaiseInternal(ctx, err)
		}
	}

	status := enum.PetitionRejected

	updateInput := dbx.UpdatePetitionInput{
		Status:  &status,
		Reply:   &reply,
		EndDate: &petition.EndDate,
	}

	if err := p.q.New().FilterID(petitionID).Update(ctx, updateInput); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Petition{}, errx.RaisePetitionNotFoundByID(ctx, err, petitionID)
		default:
			return models.Petition{}, errx.RaiseInternal(ctx, err) // Other errors
		}
	}

	//TODO add kafka event for petition rejection

	return models.Petition{
		ID:          petition.ID,
		CityID:      petition.CityID,
		CreatorID:   petition.CreatorID,
		Title:       petition.Title,
		Description: petition.Description,
		Status:      status,
		Signatures:  petition.Signatures,
		Goal:        petition.Goal,
		Reply:       reply,
		EndDate:     petition.EndDate,
		CreatedAt:   petition.CreatedAt,
		UpdatedAt:   petition.UpdatedAt,
	}, nil
}

func (p Petition) SignPetition(ctx context.Context, initiatorID, petitionID uuid.UUID) (models.PetitionSignature, error) {
	signatureID := uuid.New()
	now := time.Now().UTC()

	signature := dbx.PetitionSignature{
		ID:         signatureID,
		PetitionID: petitionID,
		UserID:     initiatorID,
		CreatedAt:  now,
	}

	if err := p.sigQ.New().Insert(ctx, signature); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			_, err = p.sigQ.New().FilterID(petitionID).FilterUserID(initiatorID).Get(ctx)
			if err == nil {
				return models.PetitionSignature{}, errx.RaisePetitionSignaturesAlreadyExists(ctx, err, petitionID, initiatorID)
			}

			_, err := p.q.New().FilterID(petitionID).Get(ctx)
			if err != nil {
				switch {
				case errors.Is(err, sql.ErrNoRows):
					return models.PetitionSignature{}, errx.RaisePetitionNotFoundByID(ctx, err, petitionID)
				}
			}
		default:
			return models.PetitionSignature{}, errx.RaiseInternal(ctx, err)
		}
	}

	return petitionSignatureModel(signature), nil
}

func (p Petition) GetSignatureByID(ctx context.Context, userID, petitionID uuid.UUID) (models.PetitionSignature, error) {
	res, err := p.sigQ.New().FilterID(petitionID).FilterUserID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.PetitionSignature{}, errx.RaisePetitionSignaturesNotFoundByPetitionIDUserID(ctx, err, petitionID, userID)
		default:
			return models.PetitionSignature{}, errx.RaiseInternal(ctx, err)
		}
	}

	return petitionSignatureModel(res), nil
}

func (p Petition) GetSignatureByUserIDAndSigID(ctx context.Context, sigID uuid.UUID) (models.PetitionSignature, error) {
	res, err := p.sigQ.New().FilterID(sigID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.PetitionSignature{}, errx.RaisePetitionSignaturesNotFoundByID(ctx, err, sigID)
		default:
			return models.PetitionSignature{}, errx.RaiseInternal(ctx, err)
		}
	}

	return petitionSignatureModel(res), nil
}

type ListPetitionsFilter struct {
	CityID    *uuid.UUID
	CreatorID *uuid.UUID
	TitleLike *string
	Rejected  *bool
	Approved  *bool
	Available *bool // Filter for available petitions (not ended)
	Expired   *bool // Filter for expired petitions
}

type ListPetitionsSort struct {
	Newest   bool // Sort by newest first
	Oldest   bool // Sort by oldest first
	MoreSign bool // Sort by more signatures first
	LessSign bool // Sort by less signatures first
}

func (p Petition) ListPetitions(
	ctx context.Context,
	filter ListPetitionsFilter,
	sort ListPetitionsSort,
	pag pagination.Request,
) ([]models.Petition, pagination.Response, error) {
	query := p.q.New()

	if filter.CityID != nil {
		query = query.FilterCityID(*filter.CityID)
	}
	if filter.CreatorID != nil {
		query = query.FilterCreatorID(*filter.CreatorID)
	}
	if filter.TitleLike != nil {
		query = query.TitleLike(*filter.TitleLike)
	}

	approved := filter.Approved != nil && *filter.Approved
	rejected := filter.Rejected != nil && *filter.Rejected
	available := filter.Available != nil && *filter.Available
	expired := filter.Expired != nil && *filter.Expired

	statuses := make([]string, 0, 3)
	if approved {
		statuses = append(statuses, enum.PetitionApproved)
	}
	if rejected {
		statuses = append(statuses, enum.PetitionRejected)
	}
	if available || expired {
		statuses = append(statuses, enum.PetitionPublished)
	}
	if len(statuses) > 0 {
		query = query.FilterStatusIn(statuses...)
	}

	now := time.Now().UTC()
	switch {
	case available && expired:
		// If both available and expired are true, we want to include all petitions that are either available or expired
	case available:
		query = query.FilterEndDate(now, true) // end_date > now
	case expired:
		query = query.FilterEndDate(now, false) // end_date < now
	}

	switch {
	case sort.MoreSign:
		query = query.OrderBySignatures(false)
	case sort.LessSign:
		query = query.OrderBySignatures(true)
	case sort.Oldest:
		query = query.OrderByCreated(true)
	default:
		query = query.OrderByCreated(false)
	}

	limit, offset := pagination.CalculateLimitOffset(pag)

	petitions, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []models.Petition{}, pagination.Response{}, nil // No petitions found, return empty slice and pagination response
		default:
			return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
		}
	}

	total, err := query.Count(ctx)
	if err != nil {
		switch {
		default:
			return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
		}
	}

	var modelsPetitions []models.Petition
	for _, p := range petitions {
		modelsPetitions = append(modelsPetitions, petitionModel(p))
	}

	return modelsPetitions, pagination.Response{
		Page:  pag.Page,
		Size:  pag.Size,
		Total: total,
	}, nil
}

type ListPetitionsSignFilter struct {
	PetitionID *uuid.UUID // Filter by specific petition ID
	UserID     *uuid.UUID // Filter by user ID who signed the petition
}

type ListPetitionsSignSort struct {
	Newest bool // Sort by newest first
	Oldest bool // Sort by oldest first
}

func (p Petition) ListSignatures(
	ctx context.Context,
	filter ListPetitionsSignFilter,
	sort ListPetitionsSignSort,
	pag pagination.Request,
) ([]models.PetitionSignature, pagination.Response, error) {
	query := p.sigQ.New()

	if filter.PetitionID != nil {
		query = query.FilterPetitionID(*filter.PetitionID)
	}
	if filter.UserID != nil {
		query = query.FilterUserID(*filter.UserID)
	}

	switch {
	case sort.Oldest:
		query = query.OrderByCreated(true)
	default:
		query = query.OrderByCreated(false)
	}

	limit, offset := pagination.CalculateLimitOffset(pag)

	signatures, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []models.PetitionSignature{}, pagination.Response{}, nil // No signatures found, return
		default:
			return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
		}
	}

	var modelsSignatures []models.PetitionSignature
	for _, sig := range signatures {
		modelsSignatures = append(modelsSignatures, petitionSignatureModel(sig))
	}

	return modelsSignatures, pagination.Response{}, errx.RaiseInternal(ctx, err)
}

func petitionModel(p dbx.Petition) models.Petition {
	return models.Petition{
		ID:          p.ID,
		CityID:      p.CityID,
		CreatorID:   p.CreatorID,
		Title:       p.Title,
		Description: p.Description,
		Status:      p.Status,
		Signatures:  p.Signatures,
		Goal:        p.Goal,
		Reply:       p.Reply,
		EndDate:     p.EndDate,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func petitionSignatureModel(sig dbx.PetitionSignature) models.PetitionSignature {
	return models.PetitionSignature{
		ID:         sig.ID,
		PetitionID: sig.PetitionID,
		UserID:     sig.UserID,
		CreatedAt:  sig.CreatedAt,
	}
}
