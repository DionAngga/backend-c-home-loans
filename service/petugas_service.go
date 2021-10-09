package service

import (
	"github.com/rysmaadit/go-template/config"
	"github.com/rysmaadit/go-template/external/jwt_client"
)

type petugasService struct {
	appConfig *config.Config
	jwtClient jwt_client.JWTClientInterface
}

type PetugasServiceInterface interface {
}

func NewPetugasService(appConfig *config.Config, jwtClient jwt_client.JWTClientInterface) *petugasService {
	return &petugasService{
		appConfig: appConfig,
		jwtClient: jwtClient,
	}
}
