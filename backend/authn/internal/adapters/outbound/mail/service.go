package mail

import (
	"github.com/faber-numeris/beholder/backend/authn/internal/infrastructure/config"
	"github.com/faber-numeris/beholder/backend/authn/internal/ports/outbound"
)

// NewService creates a new mailer service
func NewService(config config.IMailConfig) outboundport.Mailer {
	return NewMailpit(config)
}
