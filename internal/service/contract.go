package service

import (
	"github.com/nmluci/gostellar"
	"github.com/stellar-payment/sp-payment/internal/repository"
	"github.com/stellar-payment/sp-payment/pkg/dto"
)

type Service interface {
	Ping() (pingResponse dto.PublicPingResponse)
}

type service struct {
	conf       *serviceConfig
	repository repository.Repository
}

type serviceConfig struct {
}

type NewServiceParams struct {
	Repository repository.Repository
	StellarRPC *gostellar.StellarRPC
}

func NewService(params *NewServiceParams) Service {
	return &service{
		conf:       &serviceConfig{},
		repository: params.Repository,
	}
}
