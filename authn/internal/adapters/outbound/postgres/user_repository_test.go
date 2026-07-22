package postgresadapter

import (
	"context"
	"errors"
	"testing"

	"github.com/faber-numeris/beholder/authn/internal/adapters/outbound/postgres/gen"
	"github.com/faber-numeris/beholder/authn/internal/core/domain"
	outboundport "github.com/faber-numeris/beholder/authn/internal/ports/outbound"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_CreateUser(t *testing.T) {
	user := &domain.User{Email: "test@example.com"}
	passwordHash := "hashed"

	t.Run("success", func(t *testing.T) {
		q := &fakeQuerier{
			createUserFn: func(ctx context.Context, arg gen.CreateUserParams) (gen.User, error) {
				assert.Equal(t, user.Email, arg.Email)
				assert.Equal(t, passwordHash, arg.PasswordHash)
				return gen.User{ID: "123", Email: user.Email}, nil
			},
		}
		repo := NewUserRepository(q)

		res, err := repo.CreateUser(context.Background(), user, passwordHash)

		assert.NoError(t, err)
		assert.Equal(t, user.Email, res.Email)
		assert.Equal(t, "123", res.ID)
	})

	t.Run("already exists", func(t *testing.T) {
		q := &fakeQuerier{
			createUserFn: func(ctx context.Context, arg gen.CreateUserParams) (gen.User, error) {
				return gen.User{}, &pgconn.PgError{Code: pgerrcode.UniqueViolation}
			},
		}
		repo := NewUserRepository(q)

		res, err := repo.CreateUser(context.Background(), user, passwordHash)

		assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
		assert.Nil(t, res)
	})

	t.Run("error", func(t *testing.T) {
		q := &fakeQuerier{
			createUserFn: func(ctx context.Context, arg gen.CreateUserParams) (gen.User, error) {
				return gen.User{}, errors.New("db error")
			},
		}
		repo := NewUserRepository(q)

		res, err := repo.CreateUser(context.Background(), user, passwordHash)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestUserRepository_GetUserByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		q := &fakeQuerier{
			getUserFn: func(ctx context.Context, id string) (gen.User, error) {
				return gen.User{ID: "123", Email: "test@example.com"}, nil
			},
		}
		repo := NewUserRepository(q)

		res, err := repo.GetUserByID(context.Background(), "123")

		assert.NoError(t, err)
		assert.Equal(t, "123", res.ID)
	})

	t.Run("not found", func(t *testing.T) {
		q := &fakeQuerier{
			getUserFn: func(ctx context.Context, id string) (gen.User, error) {
				return gen.User{}, pgx.ErrNoRows
			},
		}
		repo := NewUserRepository(q)

		res, err := repo.GetUserByID(context.Background(), "404")

		assert.NoError(t, err)
		assert.Nil(t, res)
	})
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		q := &fakeQuerier{
			getUserByEmailFn: func(ctx context.Context, email string) (gen.User, error) {
				return gen.User{ID: "123", Email: email}, nil
			},
		}
		repo := NewUserRepository(q)

		res, err := repo.GetUserByEmail(context.Background(), "test@example.com")

		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", res.Email)
	})
}

func TestUserRepository_GetUserCredentials(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		q := &fakeQuerier{
			getUserByEmailFn: func(ctx context.Context, email string) (gen.User, error) {
				return gen.User{Email: email, PasswordHash: "hash"}, nil
			},
		}
		repo := NewUserRepository(q)

		res, err := repo.GetUserCredentials(context.Background(), "test@example.com")

		assert.NoError(t, err)
		assert.Equal(t, "hash", res.PasswordHash)
	})
}

func TestUserRepository_UpdateUser(t *testing.T) {
	user := &domain.User{ID: "123", Email: "new@example.com"}

	t.Run("success", func(t *testing.T) {
		q := &fakeQuerier{
			updateUserFn: func(ctx context.Context, arg gen.UpdateUserParams) (gen.User, error) {
				assert.Equal(t, user.ID, arg.ID)
				assert.Equal(t, user.Email, arg.Email)
				return gen.User{ID: user.ID, Email: user.Email}, nil
			},
		}
		repo := NewUserRepository(q)

		err := repo.UpdateUser(context.Background(), user)

		assert.NoError(t, err)
	})
}

func TestUserRepository_DeleteUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		q := &fakeQuerier{
			deleteUserFn: func(ctx context.Context, id string) error {
				assert.Equal(t, "123", id)
				return nil
			},
		}
		repo := NewUserRepository(q)

		err := repo.DeleteUser(context.Background(), "123")

		assert.NoError(t, err)
	})
}

func TestUserRepository_ListUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		q := &fakeQuerier{
			listUsersFn: func(ctx context.Context, arg gen.ListUsersParams) ([]gen.User, error) {
				return []gen.User{
					{ID: "1", Email: "1@ex.com"},
					{ID: "2", Email: "2@ex.com"},
				}, nil
			},
		}
		repo := NewUserRepository(q)

		res, err := repo.ListUsers(context.Background(), &outboundport.ListUsersParams{Active: true})

		assert.NoError(t, err)
		assert.Len(t, res, 2)
	})
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		q := &fakeQuerier{
			updatePasswordFn: func(ctx context.Context, arg gen.UpdatePasswordParams) error {
				assert.Equal(t, "123", arg.Userid)
				assert.Equal(t, "hash", arg.Passwordhash)
				return nil
			},
		}
		repo := NewUserRepository(q)

		err := repo.UpdatePassword(context.Background(), "123", "hash")

		assert.NoError(t, err)
	})
}

func TestUserRepository_Ping(t *testing.T) {
	assert.True(t, (&userRepository{db: &fakeQuerier{}}).Ping())
	assert.False(t, (&userRepository{}).Ping())
}
