package outboundport

import (
	"context"

	"github.com/faber-numeris/beholder/authn/internal/adapters/outbound"
	"github.com/faber-numeris/beholder/authn/internal/core/domain"
)

type Mailer interface {
	outbound.Adapter
	SendConfirmationEmail(ctx context.Context, userConfirmation domain.UserConfirmation) error
}
