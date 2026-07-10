package outboundport

import (
	"context"
	"time"

	"github.com/faber-numeris/beholder/authn/internal/adapters/outbound"
	"github.com/faber-numeris/beholder/authn/internal/core/domain"
)

type Repository interface {
	outbound.Adapter
	UserRepository
	UserConfirmationRepository
}

type UserRepository interface {
	outbound.Adapter
	CreateUser(ctx context.Context, user *domain.User, passwordHash string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserCredentials(ctx context.Context, email string) (*domain.UserCredentials, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, params *ListUsersParams) ([]*domain.User, error)
	UpdatePassword(ctx context.Context, userID string, passwordHash string) error
}

type ListUsersParams struct {
	Email             *string
	CreatedStartRange *time.Time
	CreatedEndRange   *time.Time
	Active            bool
}

type UserConfirmationRepository interface {
	outbound.Adapter
	CreateUserConfirmation(ctx context.Context, userID string, token string, expiresAt time.Time) (*domain.UserConfirmation, error)
	GetUserConfirmationByToken(ctx context.Context, token string) (string, error)
	ConfirmUserRegistration(ctx context.Context, userID string) error
	DeleteUserConfirmation(ctx context.Context, userID string) error
}
