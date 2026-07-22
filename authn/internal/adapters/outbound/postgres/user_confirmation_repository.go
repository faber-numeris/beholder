package postgresadapter

import (
	"context"
	"errors"
	"time"

	"github.com/faber-numeris/beholder/authn/internal/adapters/outbound/postgres/gen"
	"github.com/faber-numeris/beholder/authn/internal/core/domain"
	outboundport "github.com/faber-numeris/beholder/authn/internal/ports/outbound"
	"github.com/jackc/pgx/v5"
)

type userConfirmationRepository struct {
	db gen.Querier
}

func NewUserConfirmationRepository(queries gen.Querier) outboundport.UserConfirmationRepository {
	return &userConfirmationRepository{
		db: queries,
	}
}

func (r *userConfirmationRepository) Ping() bool {
	return r.db != nil
}

func (r *userConfirmationRepository) CreateUserConfirmation(ctx context.Context, userID string, token string, expiresAt time.Time) (*domain.UserConfirmation, error) {
	row, err := r.db.CreateUserConfirmation(ctx, gen.CreateUserConfirmationParams{
		Userid:    userID,
		Token:     token,
		Expiresat: expiresAt,
	})
	if err != nil {
		return nil, err
	}

	return r.toDomain(row), nil
}

func (r *userConfirmationRepository) GetUserConfirmationByToken(ctx context.Context, token string) (string, error) {
	row, err := r.db.GetUserConfirmationByToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return row.UserID, nil
}

func (r *userConfirmationRepository) ConfirmUserRegistration(ctx context.Context, userID string) error {
	return r.db.ConfirmUserRegistration(ctx, userID)
}

func (r *userConfirmationRepository) DeleteUserConfirmation(ctx context.Context, userID string) error {
	return r.db.DeleteUserConfirmation(ctx, userID)
}

func (r *userConfirmationRepository) toDomain(row gen.UserConfirmation) *domain.UserConfirmation {
	var createdAt, updatedAt time.Time
	if row.CreatedAt != nil {
		createdAt = *row.CreatedAt
	}
	if row.UpdatedAt != nil {
		updatedAt = *row.UpdatedAt
	}

	return &domain.UserConfirmation{
		ID:          row.ID,
		UserID:      row.UserID,
		Token:       row.Token,
		ExpiresAt:   row.ExpiresAt,
		ConfirmedAt: row.ConfirmedAt,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
