package model

import (
	"time"

)

type Balance struct {
	Id					int		`json:"id,string,omitempty"`
	BalanceId			string		`json:"balance_id,omitempty"`
    Account 			string 		`json:"account,omitempty"`
	Amount				int64 		`json:"amount,string,omitempty"`
    DateBalance  		time.Time 	`json:"date_balance,omitempty"`
	Description			string 		`json:"description,omitempty"`
}

type Count struct {
	Id					string		`json:"id"`
    Count 				string 		`json:"count"`
}