package util

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	DBCONNECTION string `mapstructure:"DB_CONNECTION"`

	DBHOST  string `mapstructure:"DB_HOST"`
	DBPORT  string `mapstructure:"DB_PORT"`
	DBUSER  string `mapstructure:"DB_USERNAME"`
	DBPASS  string `mapstructure:"DB_PASSWORD"`
	DBNAME  string `mapstructure:"DB_NAME"`
	SSLMODE string `mapstructure:"SSL_MODE"`

	PORTSERVER string `mapstructure:"PORT_SERVER"`
}

func LoadEnv() *Env {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("cannot read cofiguration", err)
	}

	env := &Env{}
	err = viper.Unmarshal(env)
	if err != nil {
		log.Fatal(err)
	}

	return env
}
