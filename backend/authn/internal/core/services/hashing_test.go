package services_test

import (
	"context"
	"testing"

	"github.com/alexedwards/argon2id"
	"github.com/faber-numeris/beholder/backend/authn/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashingService(t *testing.T) {
	s := services.NewHashingService()
	ctx := context.Background()

	t.Run("HashPassword success", func(t *testing.T) {
		password := []byte("password123")
		hash, err := s.HashPassword(ctx, password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)

		_, err = argon2id.ComparePasswordAndHash(string(password), hash)
		assert.NoError(t, err, "hash should be valid argon2id format and match the password")
	})

	t.Run("VerifyPassword success", func(t *testing.T) {
		password := []byte("password123")
		hash, err := s.HashPassword(ctx, password)
		require.NoError(t, err)

		match, err := s.VerifyPassword(ctx, password, hash)
		assert.NoError(t, err)
		assert.True(t, match)
	})

	t.Run("VerifyPassword wrong password", func(t *testing.T) {
		hash, err := s.HashPassword(ctx, []byte("correct-password"))
		require.NoError(t, err)

		match, err := s.VerifyPassword(ctx, []byte("wrong-password"), hash)
		assert.NoError(t, err)
		assert.False(t, match)
	})
}
