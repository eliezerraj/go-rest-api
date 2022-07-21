package model

import (
	"time"

)

type Balance struct {
	Id					string		`json:"id"`
	BalanceId			string		`json:"balance_id"`
    Account 			string 		`json:"account"`
	Amount				int32 		`json:"amount"`
    DateBalance  		time.Time 	`json:"date_balance"`
	Description			string 		`json:"description"`
}

type Count struct {
	Id					string		`json:"id"`
    Count 				string 		`json:"count"`
}