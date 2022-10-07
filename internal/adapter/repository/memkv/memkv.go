package repository

import (
	"log"
	"encoding/json"
	"sync"
	"context"

	"github.com/go-rest-api/internal/adapter/contract"
	"github.com/go-rest-api/internal/model"
	"github.com/go-rest-api/internal/error"

)

var mutex sync.Mutex

type Memkv struct {
	kv map[string][]byte
}

func NewMemKV() contract.BalanceRepositoryAdapterPort {
	return &Memkv{
		kv: map[string][]byte{},
	}
}

func (repo *Memkv) Ping() (bool, error) {
	return true, nil
}

func (repo *Memkv) ListBalance(ctx context.Context) ([]model.Balance, error) {
	log.Printf("ListBalance")

	log.Printf("####################################")
	log.Printf("- DataBase MENKV - ListBalance")
	log.Printf("####################################")

	var res []model.Balance
	for _, value := range repo.kv {
		balance := model.Balance{}
		err := json.Unmarshal(value, &balance)
		if err != nil {
			log.Printf("Erro : %v \n ", err)
			return nil, erro.ErrUnmarshal
		}
		res = append(res, balance)
	}
	return res ,nil
}

func (repo *Memkv) GetBalance(ctx context.Context, account string) (model.Balance, error) {
	log.Printf("GetBalance")

	log.Printf("####################################")
	log.Printf("- DataBase MENKV - GetBalance")
	log.Printf("####################################")

	if value, ok := repo.kv[account]; ok {
		balance := model.Balance{}
		err := json.Unmarshal(value, &balance)
		if err != nil {
			log.Printf("Erro : %v \n ", err)
			return model.Balance{}, erro.ErrUnmarshal
		}
		return balance, nil
	}
	return model.Balance{}, erro.ErrNotFound
}

func (repo *Memkv) AddBalance(ctx context.Context, balance model.Balance) (model.Balance, error) {
	log.Printf("repo-AddBalance")

	log.Printf("####################################")
	log.Printf("- DataBase MENKV - AddBalance")
	log.Printf("####################################")
	
	mutex.Lock()
	defer mutex.Unlock()

	bytes, err := json.Marshal(balance)
	if err != nil {
		log.Printf("Erro : %v \n ", err)
		return balance, erro.ErrInsert
	}
	repo.kv[balance.BalanceId] = bytes
	return balance, nil
}

func (repo *Memkv) UpdateBalance(ctx context.Context, balance model.Balance) (model.Balance, error) {
	log.Printf("repo-AddBalance")

	log.Printf("####################################")
	log.Printf("- DataBase MENKV - UpdateBalance")
	log.Printf("####################################")
	
	mutex.Lock()
	defer mutex.Unlock()

	bytes, err := json.Marshal(balance)
	if err != nil {
		log.Printf("Erro : %v \n ", err)
		return balance, erro.ErrInsert
	}
	repo.kv[balance.BalanceId] = bytes
	return balance, nil
}


func (repo *Memkv) ListBalanceById(ctx context.Context, pk string, sk string) ([]model.Balance, error) {
	return []model.Balance{}, erro.ErrFunctionNotImpl
}