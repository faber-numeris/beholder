package services

import (
	"context"

	"github.com/alexedwards/argon2id"
	inboundport "github.com/faber-numeris/beholder/backend/authn/internal/ports/outbound"
)

type hashingService struct{}

func NewHashingService() inboundport.HashingService {
	return &hashingService{}
}

func (s *hashingService) Ping() bool {
	return true
}

func (s *hashingService) HashPassword(ctx context.Context, password []byte) (string, error) {
	hash, err := argon2id.CreateHash(string(password), argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (s *hashingService) VerifyPassword(ctx context.Context, password []byte, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(string(password), hash)
	if err != nil {
		return false, err
	}
	return match, nil
}
