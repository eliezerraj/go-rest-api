package db_postgre

import (
	"log"
	"context"
	"time"
	//"strconv"
	//"encoding/json"

	"github.com/go-rest-api/internal/model"
	"github.com/go-rest-api/internal/error"
)

type BalancePostGreDBImplementacion struct {
	DatabaseHelper DatabaseHelper
}

func NewBalanceRepositoryDB(databaseHelper DatabaseHelper) BalancePostGreDBImplementacion {
	return BalancePostGreDBImplementacion{
		DatabaseHelper: databaseHelper,
	}
}

func (b BalancePostGreDBImplementacion) Ping() (bool, error) {
	log.Printf("+++++++++++++++++++++++++++++++++")
	log.Printf("- DataBase POSTGRE - Ping")
	log.Printf("+++++++++++++++++++++++++++++++++")
	_, cancel := context.WithTimeout(context.Background(), 1000)
	defer cancel()

	client, _ := b.DatabaseHelper.GetConnection(context.Background())
	err := client.Ping()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (b BalancePostGreDBImplementacion) AddBalance(ctx context.Context, balance model.Balance) (model.Balance, error) {	
	log.Printf("+++++++++++++++++++++++++++++++++")
	log.Printf("- DataBase POSTGRE - AddBalance")
	log.Printf("+++++++++++++++++++++++++++++++++")

	_, cancel := context.WithTimeout(ctx, 1000)
	defer cancel()

	client, _ := b.DatabaseHelper.GetConnection(ctx)

	stmt, err := client.Prepare(`INSERT INTO balance ( 	balance_id, 
														account, 
														amount, 
														date_balance, 
														description) 
									VALUES( $1, $2, $3, $4, $5) `)
	if err != nil {
		log.Panic(err) 
		return model.Balance{}, erro.ErrInsert
	}
	_, err = stmt.Exec(	balance.BalanceId, 
						balance.Account,
						balance.Amount,
						time.Now(),
						balance.Description)

	return balance , nil
}

func (b BalancePostGreDBImplementacion) UpdateBalance(ctx context.Context, balance model.Balance) (model.Balance, error) {	
	log.Printf("+++++++++++++++++++++++++++++++++")
	log.Printf("- DataBase POSTGRE - UpdateBalance")
	log.Printf("+++++++++++++++++++++++++++++++++")

	_, cancel := context.WithTimeout(ctx, 1000)
	defer cancel()

	client, _ := b.DatabaseHelper.GetConnection(ctx)

	stmt, err := client.Prepare(`Update balance
									set balance_id = $1, 
										account = $2, 
										amount = $3, 
										date_balance = $4, 
										description = $5
								where id = $6 `)
	if err != nil {
		log.Panic(err) 
		return model.Balance{}, erro.ErrInsert
	}
	_, err = stmt.Exec(	balance.BalanceId, 
						balance.Account,
						balance.Amount,
						time.Now(),
						balance.Description,
						balance.Id)

	return balance , nil
}

func (b BalancePostGreDBImplementacion) ListBalance(ctx context.Context) ([]model.Balance, error) {
	log.Printf("+++++++++++++++++++++++++++++++++")
	log.Printf("- DataBase POSTGRE - ListBalance")
	log.Printf("+++++++++++++++++++++++++++++++++")
	_, cancel := context.WithTimeout(ctx, 1000)
	defer cancel()

	client, _ := b.DatabaseHelper.GetConnection(ctx)

	balance := model.Balance{}
	balance_list := []model.Balance{}

	rows, err := client.Query(`SELECT id, balance_id, account, amount, date_balance, description FROM balance order by id desc`)
	defer rows.Close()
	if err != nil {
		log.Println("GetBalance ", err)
		return []model.Balance{}, err
	}

	for rows.Next() {
		err := rows.Scan( 	&balance.Id, 
							&balance.BalanceId, 
							&balance.Account, 
							&balance.Amount, 
							&balance.DateBalance , 
							&balance.Description )
		if err != nil {
            panic(err)
        }
		balance_list = append(balance_list,balance)
	}

	return balance_list , nil
}

func (b BalancePostGreDBImplementacion) GetBalance(ctx context.Context, id int) (model.Balance, error) {
	log.Printf("+++++++++++++++++++++++++++++++++")
	log.Print("- DataBase POSTGRE - GetBalance")
	log.Printf("+++++++++++++++++++++++++++++++++")
	
	_, cancel := context.WithTimeout(ctx, 1000)
	defer cancel()

	client, _ := b.DatabaseHelper.GetConnection(ctx)

	balance := model.Balance{}

	rows, err := client.Query(`SELECT id, balance_id ,account, amount, date_balance, description FROM balance WHERE id = $1`, id)
	defer rows.Close()
	if err != nil {
		log.Println("GetBalance ", err)
		return balance, err
	}

	for rows.Next() {
		err := rows.Scan( &balance.Id, &balance.BalanceId, &balance.Account, &balance.Amount, &balance.DateBalance , &balance.Description )
		if err != nil {
            panic(err)
        }
		return balance, nil
	}

	return balance , erro.ErrNotFound
}

func (b BalancePostGreDBImplementacion) ListBalanceById(ctx context.Context, pk string, sk string) ([]model.Balance, error) {
	return []model.Balance{}, erro.ErrFunctionNotImpl
}