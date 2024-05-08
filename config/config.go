package config

import "github.com/tkanos/gonfig"

type DBConfiguration struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_PORT     string
	DB_HOST     string
	DB_NAME     string
}

func GetDBConfig() DBConfiguration {
	conf := DBConfiguration{}
	gonfig.GetConf("config/config.json", &conf)
	return conf
}
