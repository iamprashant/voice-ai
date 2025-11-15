package external_emailer

import (
	external_clients "github.com/rapidaai/pkg/clients/external"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/configs"
)

func NewEmailer(config *configs.EmailerConfig, logger commons.Logger) external_clients.Emailer {
	if config == nil {
		return NewLocalEmailer(logger)
	}
	switch config.Provider() {
	case configs.SENDGRID:
		return NewSendgridEmailer(logger, config)
	case configs.SES:
		return NewSESEmailer(logger, config)
	default:
		return NewLocalEmailer(logger)
	}
}
