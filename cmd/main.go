package main

import (
	"log"
	"os"
	"time"
	"strconv"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"

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
var table_name = "balance_new"

type RDSSecrets struct {
	DB_USER    	string
	DB_PASSWORD string 
}

func init(){
	log.Printf("==========================")
	log.Printf("Init")
	getEnvfromFile()
	getEnv()
	if (app.AwsEnv.IcUseSsm){
		getSSMSecrets()
	} 
	log.Printf("==========================")
}

func getEnvfromFile() {
	log.Printf("-------------------")
	log.Printf("Loading variables from application.yaml")
	app, err = config.Configuration()
	if err != nil {
		log.Print("FATAL ERROR load application.yaml", err)
		panic(err)	
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
	if os.Getenv("DB_NAME") !=  "" {	
		app.DatabaseRDS.DatabaseName = os.Getenv("DB_NAME")
	}
	if os.Getenv("DB_SCHEMA") !=  "" {	
		app.DatabaseRDS.Schema = os.Getenv("DB_SCHEMA")
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
	if os.Getenv("AWS_SSM_SECRET_ID") !=  "" {	
		app.AwsEnv.AwsSsmSecretId = os.Getenv("AWS_SSM_SECRET_ID")
	}
	if os.Getenv("IC_USE_SSM") !=  "" {	
		boolVar, _ := strconv.ParseBool(os.Getenv("IC_USE_SSM"))
		app.AwsEnv.IcUseSsm = boolVar
	}
	config.PodInfo(&app)
	log.Println("Variables", app)
}

func getSSMSecrets(){
	log.Printf("-------------------")
	log.Printf("get Secrets from SSM")

	svc := secretsmanager.New(session.New(&aws.Config {
		Region: &app.AwsEnv.Aws_region,
	}))

	ssm_secrets := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(app.AwsEnv.AwsSsmSecretId),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(ssm_secrets)
	if err != nil {
		log.Print("FATAL ERROR get secret SSM ", err)
		panic(err)	
	}

	var rds_secrets RDSSecrets

	var secretString string
	if result.SecretString != nil {
		secretString = *result.SecretString
		err := json.Unmarshal([]byte(secretString) , &rds_secrets)
		if err != nil {
			log.Println("FATAL ERROR conversion Credentials JSON ", err)
			panic(err)	
		}
		app.DatabaseRDS.User = rds_secrets.DB_USER
		app.DatabaseRDS.Password = rds_secrets.DB_PASSWORD
		log.Print("databaseRDS " , &app.DatabaseRDS)
	}
}

func main(){
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting Http Server 1.0" )

	var repo_db contract.BalanceRepositoryAdapterPort

	if (app.Setup.DatabaseType == "POSTGRES"){
		log.Printf("==> Database POSTGRES")

		count := 1
		var config_postgre db_postgre.DatabaseHelper

		for {
			config_postgre, err = db_postgre.NewDatabaseHelper(app.DatabaseRDS)
			if err != nil {
				if count < 3 {
					log.Printf("Erro na abertura do Database try %v ... err: %v",count ,err)
				} else {
					log.Printf("ERRO FATAL na abertura do Database aborting try %v ... err: %v", count ,err)
					panic(err)	
				}
				time.Sleep(3 * time.Second)
				count = count + 1
				continue
			}
			break
		}

		repo_db = db_postgre.NewBalanceRepositoryDB(config_postgre)
	}else if (app.Setup.DatabaseType == "DYNAMO"){
		log.Printf("==> Database DYNAMO")
		repo_db, err = dynamodb.NewBalanceRepositoryDynamoDB(table_name,app)
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