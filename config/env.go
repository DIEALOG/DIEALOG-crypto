package config

import (
	"log"

	"github.com/spf13/viper"
)

var E Env

type Env struct {
	AppEnv            string `mapstructure:"APP_ENV"`
	ServerAddress     string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout    int    `mapstructure:"CONTEXT_TIMEOUT"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPass            string `mapstructure:"DB_PASS"`
	DBName            string `mapstructure:"DB_NAME"`
	InternalSecretKey string `mapstructure:"INTERNAL_SECRET_KEY"` // Public API <-> Internal Crypto Service 간 시크릿
	KMSKeyARN         string `mapstructure:"KMS_KEY_ARN"`         // AWS KMS Key ARN
	MockKMSMasterKey  string `mapstructure:"MOCK_KMS_MASTER_KEY"`
}

func NewEnv() {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .E : ", err)
	}

	err = viper.Unmarshal(&E)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if E.AppEnv == "development" {
		log.Println("The App is running in development E")
	}
}
