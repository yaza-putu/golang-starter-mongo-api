package config

import "github.com/spf13/viper"

type db struct {
	Host        string
	User        string
	Name        string
	Password    string
	Port        int
	AutoMigrate bool
}

func DB() db {
	return db{
		Host:        viper.GetString("mongo_host"),
		User:        viper.GetString("mongo_user"),
		Name:        viper.GetString("mongo_database"),
		Password:    viper.GetString("mongo_password"),
		Port:        viper.GetInt("mongo_port"),
		AutoMigrate: viper.GetBool("mongo_auto_migrate"),
	}
}
