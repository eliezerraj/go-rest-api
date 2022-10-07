# go-mock-api

How to setup

Dynamo - tablename - balance {balance_id (pk) account (sk)}

CREATE TABLE public.balance (
    id              int serial, 
	balance_id      int4 NULL,
	account         varchar(200) NULL,
	amount          float8 NULL,
	date_balance    timestamptz NULL,
	description     varchar(200) NULL
);

1) resources/application.yml

1.1 RDS 
setup:
  databaseType: "rds" or "dynamo" or "menkv"
  responseTime: 8900
  responseStatusCode: 200
  isRandomTime: 60
  count: 0

1.2 Dynamo

awsenv:
  aws_region: ""
  aws_access_id: ""
  aws_access_secret: ""

1.3 MenKv

GET http://localhost:8900/balance/list
GET http://localhost:8900/balance/6
GET http://localhost:8900/info
GET http://localhost:8900/count/3
POST http://localhost:8900/stress/cpu
    {
        "count":2000
    }

GET http://localhost:8900/balance/list_by_id/1&sk=acc
POST http://localhost:8900/balance/save
    {
        "id": "77",
        "balance_id": "1",
        "account": "acc-3",
        "amount": 10,
        "date_balance": "2020-01-01T00:00:00Z",
        "description": "tste"
    }

POST http://localhost:8900/setup
    {
        "response_time" : 1,
        "response_status_code" : 500,
        "is_random_time": true,
        "count":0,
        "liveness": true,
        "readiness": true
    }