package mail

import (
	"github.com/faber-numeris/beholder/authn/internal/infrastructure/config"
)

// NewService creates a new mailer service
func NewService(config config.IMailConfig) outboundport.Mailer {
	return NewMailpit(config)
}
