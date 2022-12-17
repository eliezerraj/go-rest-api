package main

import (
	"log"
	"os"
	"time"
	"strconv"

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
	log.Printf("==========================")
	log.Printf("Init")
	getEnvfromFile()
	getEnv()
	log.Printf("==========================")
}

func getEnvfromFile() {
	log.Printf("-------------------")
	log.Printf("Loading variables from application.yaml")
	app, err = config.Configuration()
	if err != nil {
		log.Print("FATAL ERROR load application.yaml", err)
		os.Exit(3)
	}
	config.PodInfo(&app)
	log.Println("Variables", app)
}

func getEnv() {
	log.Printf("-------------------")
	log.Printf("Overiding variables from enviroment")
	if os.Getenv("PORT") !=  "" {
		intVar, _ := strconv.Atoi(os.Getenv("PORT"))
		app.Server.Port = intVar
	}
	if os.Getenv("API_NAME") !=  "" {
		app.App.Name = os.Getenv("API_NAME")
	}
	if os.Getenv("API_VERSION") !=  "" {
		app.App.Version = os.Getenv("API_VERSION")
	}
	if os.Getenv("DATABASE_TYPE") !=  "" {
		app.Setup.DatabaseType = os.Getenv("DATABASE_TYPE")
	}
	if os.Getenv("DB_HOST") !=  "" {
		app.DatabaseRDS.Host = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_PORT") !=  "" {
		app.DatabaseRDS.Port = os.Getenv("DB_PORT")
	}
	if os.Getenv("DB_USER") !=  "" {
		app.DatabaseRDS.User = os.Getenv("DB_USER")
	}
	if os.Getenv("DB_PASSWORD") !=  "" {	
		app.DatabaseRDS.Password = os.Getenv("DB_PASSWORD")
	}

	if os.Getenv("AWS_REGION") !=  "" {	
		app.AwsEnv.Aws_region = os.Getenv("AWS_REGION")
	}
	if os.Getenv("AWS_ACCESS_KEY") !=  "" {	
		app.AwsEnv.Aws_access_id = os.Getenv("AWS_ACCESS_KEY")
	}
	if os.Getenv("AWS_ACCESS_SECRET") !=  "" {	
		app.AwsEnv.Aws_access_secret = os.Getenv("AWS_ACCESS_SECRET")
	}

	config.PodInfo(&app)
	log.Println("Variables", app)
}

func main(){
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting Http Server 1.0")

	var repo_db contract.BalanceRepositoryAdapterPort

	if (app.Setup.DatabaseType == "POSTGRES"){
		log.Printf("==> Database POSTGRES")
		config_postgre, err := db_postgre.NewDatabaseHelper(app.DatabaseRDS)
		if err != nil {
			log.Print("Erro na abertura do Database", err)
			panic(err)
		}
		repo_db = db_postgre.NewBalanceRepositoryDB(config_postgre)
	}else if (app.Setup.DatabaseType == "DYNAMO"){
		log.Printf("==> Database DYNAMO")
		repo_db, err = dynamodb.NewBalanceRepositoryDynamoDB("balance",app)
		if err != nil {
			log.Print("Erro na abertura do Database", err)
			panic(err)
		}
	}else {
		log.Printf("==> Database MENKV")
		repo_db = repository.NewMemKV()
	}

	service_balance := controller.NewBalanceService(repo_db)
	service_metrics := controller.NewMetricsService()
	handler_balance	:= http.NewBalanceHttpAdapter(service_metrics, service_balance, repo_db)

	httpServer 		:= http.NewHttpServer(time.Now(), app)
	httpServer.StartHttpServer(handler_balance)
}