package mail

import (
	"github.com/rs/zerolog"
)

type MailtrapConfig struct {
}

type MailtrapProvider struct {
	config *MailConfig
	logger *zerolog.Logger
}

func NewMailtrapProvider(config *MailConfig) (EmailProviderService, error) {

	return &MailtrapProvider{
		config: &MailtrapConfig{},
		logger: config.Logger,
	}, nil
}
