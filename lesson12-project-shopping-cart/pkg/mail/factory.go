package mail

import (
	"fmt"
	"user-management-api/internal/utils"
)

type ProviderType string

const (
	ProviderMailtrap ProviderType = "mailtrap"
)

type ProviderFactory interface {
	CreateProvider(config *MailConfig) (EmailProviderService, error)
}

type MailtrapProviderFactory struct{}

func (f *MailtrapProviderFactory) CreateProvider(config *MailConfig) (EmailProviderService, error) {
	return NewMailtrapProvider(config)
}

func NewProviderFactory(providerType ProviderType) (ProviderFactory, error) {
	switch providerType {
	case ProviderMailtrap:
		return &MailtrapProviderFactory{}, nil
	default:
		return nil, utils.NewError(fmt.Sprintf("Unsupported provider type: %s", utils.ErrorCode(providerType)), utils.ErrCodeInternal)
	}

}
