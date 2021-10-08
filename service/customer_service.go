package service

import (
	"github.com/rysmaadit/go-template/config"
	"github.com/rysmaadit/go-template/external/jwt_client"
)

type customerService struct {
	appConfig *config.Config
	jwtClient jwt_client.JWTClientInterface
}

type CustomerServiceInterface interface {
}

func NewCustomerService(appConfig *config.Config, jwtClient jwt_client.JWTClientInterface) *customerService {
	return &customerService{
		appConfig: appConfig,
		jwtClient: jwtClient,
	}
}
