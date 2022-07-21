package controller

import (
	"log"
	"context"

	"github.com/go-rest-api/internal/model"
	"github.com/go-rest-api/internal/adapter/contract"
)

type service struct {
	repository 	contract.BalanceRepositoryAdapterPort
	ctx 		context.Context
}

func NewBalanceService(repository contract.BalanceRepositoryAdapterPort) *service {
	return &service{
		repository: repository,
		ctx: context.Background(),
	}
}

func (p *service) AddBalance(balance model.Balance) (model.Balance, error) {
	log.Printf("AddBalance")
	
	log.Printf("--------------------------------------")
	log.Printf("- AddBalance Doing Business Rules !!!!")
	log.Printf("--------------------------------------")

	_, err := p.repository.AddBalance(p.ctx, balance)
	if err != nil {
		return balance, err
	}

	return balance, nil
}

func (p *service) ListBalance() ([]model.Balance, error) {
	log.Printf("ListBalance")
	
	log.Printf("--------------------------------------")
	log.Printf("- List Doing Business Rules !!!!")
	log.Printf("--------------------------------------")

	res, err := p.repository.ListBalance(p.ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *service) ListBalanceById(pk string, sk string) ([]model.Balance, error) {
	log.Printf("ListBalanceById")
	
	log.Printf("--------------------------------------")
	log.Printf("- ListBalanceById Doing Business Rules !!!!")
	log.Printf("--------------------------------------")

	res, err := p.repository.ListBalanceById(p.ctx, pk, sk)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *service) GetBalance(account string) (model.Balance, error) {
	log.Printf("GetBalance")

	log.Printf("--------------------------------------")
	log.Printf("- GetBalance Doing Business Rules !!!!")
	log.Printf("--------------------------------------")

	res, err := p.repository.GetBalance(p.ctx, account)
	if err != nil {
		return model.Balance{}, err
	}
	
	return res, nil
}
