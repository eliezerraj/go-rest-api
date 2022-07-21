package config

import (
	"log"
	"github.com/spf13/viper"
	"github.com/go-rest-api/internal/model"
)

var (
	Application model.ManagerInfo
)

const (
	_applicationFileName = "application"
	_extension           = "yml"
	_resourcePath        = "../resource/"
)

func Configuration() (model.ManagerInfo, error){

	var app model.ManagerInfo

	viper.SetConfigName(_applicationFileName)
	viper.SetConfigType(_extension)
	viper.AddConfigPath(_resourcePath)
	viper.ReadInConfig()

	errUnmarshal := viper.Unmarshal(&app)
	if errUnmarshal != nil {
		log.Printf("Parse error for application structure", errUnmarshal)
		return app, errUnmarshal
	}

	return app, nil
}