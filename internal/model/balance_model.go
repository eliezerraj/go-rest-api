package model

import (
	"time"

)

type Balance struct {
	Id					int			`json:"id,int,omitempty"`
	BalanceId			string		`json:"balance_id,omitempty"`
    Account 			string 		`json:"account,omitempty"`
	Amount				float32		`json:"amount,float,omitempty"`
    DateBalance  		time.Time 	`json:"date_balance,omitempty"`
	Description			string 		`json:"description,omitempty"`
}

type Count struct {
	Id					string		`json:"id"`
    Count 				string 		`json:"count"`
}