package outboundport

import (
	"context"

	"github.com/faber-numeris/beholder/backend/authn/internal/adapters/outbound"
)

type HashingService interface {
	outbound.Adapter
	HashPassword(ctx context.Context, password []byte) (string, error)
	VerifyPassword(ctx context.Context, password []byte, hash string) (bool, error)
}
