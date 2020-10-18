package config

import (
	"github.com/spf13/viper"
	"log"
)

var (
	Conf = newConfig()
)

type Config struct {
	Db    ConfDB    `mapstructure:"database"`
	Web   ConfWeb   `mapstructure:"web"`
}

type ConfDB struct {
	Postgres ConfPostgres `mapstructure:"postgres"`
}

type ConfPostgres struct {
	DriverName string `mapstructure:"driver_name"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	DbName     string `mapstructure:"db_name"`
	SslMode    string `mapstructure:"ssl_mode"`
	Host       string `mapstructure:"host"`
	MaxConn    string `mapstructure:"max_conn"`
}

type ConfWeb struct {
	Server ConfServer `mapstructure:"server"`
}


type ConfServer struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
}


func newConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs/yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	conf := new(Config)

	er := viper.Unmarshal(conf)
	if er != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	return conf
}
