package clients

import (
	"field-service/clients/config"
	clients "field-service/clients/user"
	config2 "field-service/config"
)

type ClientRegistry struct{}

type IClientRegistry interface {
	GetUser() clients.IUserClient
}

func NewClientRegistry() IClientRegistry {
	return &ClientRegistry{}
}

func (c *ClientRegistry) GetUser() clients.IUserClient {
	return clients.NewUserClient(
		config.NewClientConfig(
			config.WithBaseURL(config2.Config.InternalService.User.Host),  // use localhost:8001 for local development | use user-service:8001 for docker
			config.WithSignatureKey(config2.Config.InternalService.User.SignatureKey),
		))
}
