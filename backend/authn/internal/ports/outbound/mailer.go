package outboundport

import (
	"context"

	"github.com/faber-numeris/beholder/backend/authn/internal/adapters/outbound"
	"github.com/faber-numeris/beholder/backend/authn/internal/core/domain"
)

type Mailer interface {
	outbound.Adapter
	SendConfirmationEmail(ctx context.Context, userConfirmation domain.UserConfirmation) error
}
