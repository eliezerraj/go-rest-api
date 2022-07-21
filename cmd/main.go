package main

import (
	"log"
	"os"
	"time"

	"github.com/go-rest-api/internal/config"
	"github.com/go-rest-api/internal/model"
	"github.com/go-rest-api/internal/controller"
	"github.com/go-rest-api/internal/handler/http"
	"github.com/go-rest-api/internal/adapter/contract"
	"github.com/go-rest-api/internal/adapter/repository/memkv"
	"github.com/go-rest-api/internal/adapter/repository/dynamodb"
	"github.com/go-rest-api/internal/adapter/repository/db_postgre"
)

var app model.ManagerInfo
var err error

func init(){
	log.Printf("-------------------")
	log.Printf("Loading init data")
	app, err = config.Configuration()
	if err != nil {
		log.Print("FATAL ERROR load application.yaml", err)
		os.Exit(3)
	}
	config.PodInfo(&app)
	log.Printf("-------------------")
}

func main(){
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting Http Server 1.0")
	log.Print("pod :", app)

	var repo_db contract.BalanceRepositoryAdapterPort

	if (app.Setup.DatabaseType == "postgre"){
		config_postgre, err := db_postgre.NewDatabaseHelper(app.DatabaseRDS)
		if err != nil {
			log.Print("Erro na abertura do Database", err)
			panic(err)
		}
		repo_db 	= db_postgre.NewBalanceRepositoryDB(config_postgre)
	}else if (app.Setup.DatabaseType == "dynamodb"){
		repo_db, err = dynamodb.NewBalanceRepositoryDynamoDB("balance",app)
		if err != nil {
			log.Print("Erro na abertura do Database", err)
			panic(err)
		}
	}else {
		repo_db = repository.NewMemKV()
	}

	service_balance := controller.NewBalanceService(repo_db)
	service_metrics := controller.NewMetricsService()
	handler_balance	:= http.NewBalanceHttpAdapter(service_metrics, service_balance, repo_db)

	httpServer 		:= http.NewHttpServer(time.Now(), app)
	httpServer.StartHttpServer(handler_balance)
}