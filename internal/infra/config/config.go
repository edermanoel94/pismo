package config

import "strings"

import (
	"github.com/spf13/viper"
	"log"
)

var config = viper.New()

func Init() *viper.Viper {
	config.AddConfigPath(".")
	config.AddConfigPath("/app/internal/infra/config/")
	config.SetConfigName("configuration")
	config.SetConfigType("yaml")

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			log.Fatalf("Error reading config file: %s", err)
		} else {
			log.Printf("Load default configs")
		}
	}

	setConfigDefaults()
	return config
}

func setConfigDefaults() {

	config.SetDefault("server.addr", "0.0.0.0:8080")
	config.SetDefault("server.timeout.read-seconds", "15")
	config.SetDefault("server.timeout.write-seconds", "20")
	config.SetDefault("server.debug", true)

	config.SetDefault("db.host", "localhost")
	config.SetDefault("db.user", "postgres")
	config.SetDefault("db.password", "pismo")
	config.SetDefault("db.name", "pismo")

	config.SetDefault("operation_types.COMPRA_A_VISTA", "-")
	config.SetDefault("operation_types.COMPRA_PARCELADA", "-")
	config.SetDefault("operation_types.SAQUE", "-")
	config.SetDefault("operation_types.PAGAMENTO", "+")
	config.SetDefault("operation_types.LIMITE_DE_CREDITO", "+")

	config.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	config.AutomaticEnv()
}

func Config() *viper.Viper {
	return config
}
