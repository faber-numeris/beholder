package postgresadapter

import (
	"context"
	"errors"

	"github.com/faber-numeris/beholder/authn/internal/adapters/outbound/postgres/gen"
	domain2 "github.com/faber-numeris/beholder/authn/internal/core/domain"
	outboundport "github.com/faber-numeris/beholder/authn/internal/ports/outbound"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type userRepository struct {
	db gen.Querier
}

func NewUserRepository(queries gen.Querier) outboundport.UserRepository {
	return &userRepository{
		db: queries,
	}
}

func (r *userRepository) Ping() bool {
	// TODO: request the actual database to verify if it is correctly working
	return r.db != nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain2.User, passwordHash string) (*domain2.User, error) {
	var firstName, lastName, locale, timezone string
	if user.Profile != nil {
		firstName = user.Profile.FirstName
		lastName = user.Profile.LastName
		locale = user.Profile.Locale
		timezone = user.Profile.Timezone
	}

	row, err := r.db.CreateUser(ctx, gen.CreateUserParams{
		Email:        user.Email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		Locale:       locale,
		Timezone:     timezone,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return nil, domain2.ErrUserAlreadyExists
		}
		return nil, err
	}

	return r.toDomain(row), nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*domain2.User, error) {
	row, err := r.db.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(row), nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain2.User, error) {
	row, err := r.db.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(row), nil
}

func (r *userRepository) GetUserCredentials(ctx context.Context, email string) (*domain2.UserCredentials, error) {
	row, err := r.db.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &domain2.UserCredentials{
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
	}, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *domain2.User) error {
	var firstName, lastName, locale, timezone string
	if user.Profile != nil {
		firstName = user.Profile.FirstName
		lastName = user.Profile.LastName
		locale = user.Profile.Locale
		timezone = user.Profile.Timezone
	}

	_, err := r.db.UpdateUser(ctx, gen.UpdateUserParams{
		Email:     user.Email,
		FirstName: firstName,
		LastName:  lastName,
		Locale:    locale,
		Timezone:  timezone,
		ID:        user.ID,
	})
	return err
}

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	return r.db.DeleteUser(ctx, id)
}

func (r *userRepository) ListUsers(ctx context.Context, params *outboundport.ListUsersParams) ([]*domain2.User, error) {
	var email pgtype.Text
	if params.Email != nil {
		email = pgtype.Text{String: *params.Email, Valid: true}
	}

	rows, err := r.db.ListUsers(ctx, gen.ListUsersParams{
		Email:             email,
		CreatedStartRange: params.CreatedStartRange,
		CreatedEndRange:   params.CreatedEndRange,
		Active:            params.Active,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*domain2.User, len(rows))
	for i, row := range rows {
		result[i] = r.toDomain(row)
	}

	return result, nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID string, passwordHash string) error {
	return r.db.UpdatePassword(ctx, gen.UpdatePasswordParams{
		Passwordhash: passwordHash,
		Userid:       userID,
	})
}

func (r *userRepository) toDomain(row gen.User) *domain2.User {
	// TODO: use goverter here
	return &domain2.User{
		ID:    row.ID,
		Email: row.Email,
		Profile: &domain2.UserProfile{
			FirstName: row.FirstName,
			LastName:  row.LastName,
			Locale:    row.Locale,
			Timezone:  row.Timezone,
		},
	}
}
