package postgresadapter

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/faber-numeris/beholder/authn/internal/adapters/outbound/postgres/gen"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestUserConfirmationRepository_CreateUserConfirmation(t *testing.T) {
	userID := "user-123"
	token := "token-123"
	expiresAt := time.Now().Add(time.Hour)

	t.Run("success", func(t *testing.T) {
		q := &fakeQuerier{
			createUserConfirmationFn: func(ctx context.Context, arg gen.CreateUserConfirmationParams) (gen.UserConfirmation, error) {
				assert.Equal(t, userID, arg.Userid)
				assert.Equal(t, token, arg.Token)
				return gen.UserConfirmation{ID: "conf-123", UserID: userID, Token: token, ExpiresAt: expiresAt}, nil
			},
		}
		repo := NewUserConfirmationRepository(q)

		res, err := repo.CreateUserConfirmation(context.Background(), userID, token, expiresAt)

		assert.NoError(t, err)
		assert.Equal(t, userID, res.UserID)
		assert.Equal(t, token, res.Token)
	})
}

func TestUserConfirmationRepository_GetUserConfirmationByToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		q := &fakeQuerier{
			getUserConfirmationByTokenFn: func(ctx context.Context, token string) (gen.UserConfirmation, error) {
				return gen.UserConfirmation{UserID: "user-123"}, nil
			},
		}
		repo := NewUserConfirmationRepository(q)

		userID, err := repo.GetUserConfirmationByToken(context.Background(), "valid-token")

		assert.NoError(t, err)
		assert.Equal(t, "user-123", userID)
	})

	t.Run("not found", func(t *testing.T) {
		q := &fakeQuerier{
			getUserConfirmationByTokenFn: func(ctx context.Context, token string) (gen.UserConfirmation, error) {
				return gen.UserConfirmation{}, pgx.ErrNoRows
			},
		}
		repo := NewUserConfirmationRepository(q)

		userID, err := repo.GetUserConfirmationByToken(context.Background(), "invalid")

		assert.NoError(t, err)
		assert.Equal(t, "", userID)
	})

	t.Run("error", func(t *testing.T) {
		q := &fakeQuerier{
			getUserConfirmationByTokenFn: func(ctx context.Context, token string) (gen.UserConfirmation, error) {
				return gen.UserConfirmation{}, errors.New("db error")
			},
		}
		repo := NewUserConfirmationRepository(q)

		userID, err := repo.GetUserConfirmationByToken(context.Background(), "any")

		assert.Error(t, err)
		assert.Equal(t, "", userID)
	})
}

func TestUserConfirmationRepository_ConfirmUserRegistration(t *testing.T) {
	userID := "user-123"

	t.Run("success", func(t *testing.T) {
		q := &fakeQuerier{
			confirmUserRegistrationFn: func(ctx context.Context, uid string) error {
				assert.Equal(t, userID, uid)
				return nil
			},
		}
		repo := NewUserConfirmationRepository(q)

		err := repo.ConfirmUserRegistration(context.Background(), userID)

		assert.NoError(t, err)
	})
}

func TestUserConfirmationRepository_DeleteUserConfirmation(t *testing.T) {
	userID := "user-123"

	t.Run("success", func(t *testing.T) {
		q := &fakeQuerier{
			deleteUserConfirmationFn: func(ctx context.Context, uid string) error {
				assert.Equal(t, userID, uid)
				return nil
			},
		}
		repo := NewUserConfirmationRepository(q)

		err := repo.DeleteUserConfirmation(context.Background(), userID)

		assert.NoError(t, err)
	})
}

func TestUserConfirmationRepository_Ping(t *testing.T) {
	assert.True(t, (&userConfirmationRepository{db: &fakeQuerier{}}).Ping())
	assert.False(t, (&userConfirmationRepository{}).Ping())
}
