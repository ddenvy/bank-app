package service

import (
	"bank-app/internal/config"
	"bank-app/internal/repository"
)

type Services struct {
	Users     UserService
	Accounts  AccountService
	Cards     CardService
	Credits   CreditService
	Transfers TransferService
	Analytics AnalyticsService
}

func NewServices(repos *repository.Repositories, cfg *config.Config) *Services {
	return &Services{
		Users:     NewUserService(repos.Users),
		Accounts:  NewAccountService(repos.Accounts),
		Cards:     NewCardService(repos.Cards),
		Credits:   NewCreditService(repos.Credits, repos.Accounts, cfg),
		Transfers: NewTransferService(repos.Transfers, repos.Accounts),
		Analytics: NewAnalyticsService(repos.Analytics),
	}
}

type Dependencies struct {
	Repos  *repository.Repositories
	Config *config.Config
}
